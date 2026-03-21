package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	dbmodel "practice_IstraNet/pkg/DB_model"
	"practice_IstraNet/pkg/dto"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type AdminPostgres struct {
	db *sqlx.DB
}

func NewAdminPostgres(db *sqlx.DB) *AdminPostgres {
	return &AdminPostgres{db: db}
}

func (r *AdminPostgres) AssignApplication(ctx context.Context, adminID int64, in dto.AdminAssignApplicationRequest) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	state, err := lockApplicationState(ctx, tx, in.ID)
	if err != nil {
		return err
	}
	if state.IsFinal {
		return ErrFinalApplication
	}

	exists, err := isOperatorUserTx(ctx, tx, in.OperatorID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrOperatorNotFound
	}

	statusID, err := currentStatusIDForAssign(ctx, tx, state.StatusCode)
	if err != nil {
		return err
	}
	if state.StatusCode == "new" {
		statusID, err = getStatusIDByCodeTx(ctx, tx, "in_progress")
		if err != nil {
			return err
		}
	}

	if _, err = tx.ExecContext(ctx, `
		UPDATE applications
		SET assigned_to = $1,
		    status_id = $2,
		    updated_at = now()
		WHERE id = $3
	`, in.OperatorID, statusID, in.ID); err != nil {
		return err
	}

	if state.AssignedTo == nil || *state.AssignedTo != in.OperatorID {
		var oldValue *string
		if state.AssignedTo != nil {
			oldOperator := fmt.Sprintf("%d", *state.AssignedTo)
			oldValue = &oldOperator
		}
		newOperator := fmt.Sprintf("%d", in.OperatorID)
		if err = insertHistoryTx(ctx, tx, in.ID, adminID, "assign", strPtr("assigned_to"), oldValue, &newOperator); err != nil {
			return err
		}
	}

	if state.StatusCode == "new" {
		if err = insertHistoryTx(ctx, tx, in.ID, adminID, "status_change", strPtr("status_id"), strPtr("new"), strPtr("in_progress")); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *AdminPostgres) ChangeApplicationStatusByAdmin(ctx context.Context, adminID int64, in dto.AdminChangeApplicationStatusRequest) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	state, err := lockApplicationState(ctx, tx, in.ID)
	if err != nil {
		return err
	}

	statusID, err := getStatusIDByCodeTx(ctx, tx, in.StatusCode)
	if err != nil {
		if errors.Is(err, ErrStatusNotFound) {
			return ErrInvalidStatusCode
		}
		return err
	}

	if _, err = tx.ExecContext(ctx, `
		UPDATE applications
		SET status_id = $1,
		    updated_at = now(),
		    closed_at = CASE
		        WHEN $2 = 'closed' THEN now()
		        WHEN $2 <> 'closed' THEN NULL
		        ELSE closed_at
		    END
		WHERE id = $3
	`, statusID, in.StatusCode, in.ID); err != nil {
		return err
	}

	if state.StatusCode != in.StatusCode {
		if err = insertHistoryTx(ctx, tx, in.ID, adminID, "status_change", strPtr("status_id"), strPtr(state.StatusCode), &in.StatusCode); err != nil {
			return err
		}
	}

	if in.Comment != nil {
		if _, err = tx.ExecContext(ctx, `
			INSERT INTO application_comments(application_id, author_id, body)
			VALUES ($1, $2, $3)
		`, in.ID, adminID, *in.Comment); err != nil {
			return err
		}
		if err = insertHistoryTx(ctx, tx, in.ID, adminID, "comment", nil, nil, in.Comment); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *AdminPostgres) DeleteApplicationByAdmin(ctx context.Context, appID int64) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM applications WHERE id=$1`, appID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrAppNotFound
	}
	return nil
}

func (r *AdminPostgres) GetUsers(ctx context.Context) ([]dbmodel.AdminUserDB, error) {
	var items []dbmodel.AdminUserDB

	q := `
		SELECT
			u.id,
			u.email,
			u.full_name,
			r.code AS role_code,
			r.name AS role_name,
			u.is_active,
			u.created_at
		FROM users u
		JOIN user_roles ur ON ur.user_id = u.id
		JOIN roles r ON r.id = ur.role_id
		ORDER BY u.created_at DESC, u.id DESC
	`
	if err := r.db.SelectContext(ctx, &items, q); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *AdminPostgres) GetUserByIDForAdmin(ctx context.Context, userID int64) (dbmodel.AdminUserDB, error) {
	var item dbmodel.AdminUserDB

	q := `
		SELECT
			u.id,
			u.email,
			u.full_name,
			r.code AS role_code,
			r.name AS role_name,
			u.is_active,
			u.created_at
		FROM users u
		JOIN user_roles ur ON ur.user_id = u.id
		JOIN roles r ON r.id = ur.role_id
		WHERE u.id = $1
	`
	if err := r.db.GetContext(ctx, &item, q, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dbmodel.AdminUserDB{}, ErrUserNotFound
		}
		return dbmodel.AdminUserDB{}, err
	}
	return item, nil
}

func (r *AdminPostgres) ChangeUserRole(ctx context.Context, adminID int64, in dto.AdminChangeUserRoleRequest) error {
	if adminID == in.UserID && in.RoleCode != "admin" {
		return ErrSelfActionForbidden
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	var userExists bool
	if err = tx.GetContext(ctx, &userExists, `SELECT EXISTS(SELECT 1 FROM users WHERE id=$1)`, in.UserID); err != nil {
		return err
	}
	if !userExists {
		return ErrUserNotFound
	}

	var roleID int64
	if err = tx.GetContext(ctx, &roleID, `SELECT id FROM roles WHERE code=$1`, in.RoleCode); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrRoleNotFound
		}
		return err
	}

	result, err := tx.ExecContext(ctx, `UPDATE user_roles SET role_id=$1 WHERE user_id=$2`, roleID, in.UserID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		if _, err = tx.ExecContext(ctx, `INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2)`, in.UserID, roleID); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *AdminPostgres) DeleteUserByAdmin(ctx context.Context, adminID, userID int64) error {
	if adminID == userID {
		return ErrSelfActionForbidden
	}

	res, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE id=$1`, userID)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23503" {
			return ErrUserHasRelations
		}
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *AdminPostgres) CreateStatus(ctx context.Context, in dto.AdminCreateStatusRequest) (dbmodel.ApplicationStatusDB, error) {
	var item dbmodel.ApplicationStatusDB

	q := `
		INSERT INTO application_statuses(code, name, is_final)
		VALUES ($1, $2, $3)
		RETURNING id, code, name, is_final
	`
	if err := r.db.GetContext(ctx, &item, q, in.Code, in.Name, in.IsFinal); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return dbmodel.ApplicationStatusDB{}, ErrStatusExists
		}
		return dbmodel.ApplicationStatusDB{}, err
	}
	return item, nil
}

func (r *AdminPostgres) UpdateStatus(ctx context.Context, in dto.AdminUpdateStatusRequest) (dbmodel.ApplicationStatusDB, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return dbmodel.ApplicationStatusDB{}, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	var current dbmodel.ApplicationStatusDB
	if err = tx.GetContext(ctx, &current, `
		SELECT id, code, name, is_final
		FROM application_statuses
		WHERE id = $1
		FOR UPDATE
	`, in.ID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dbmodel.ApplicationStatusDB{}, ErrStatusNotFound
		}
		return dbmodel.ApplicationStatusDB{}, err
	}

	if isProtectedStatusCode(current.Code) {
		return dbmodel.ApplicationStatusDB{}, ErrStatusProtected
	}

	if err = tx.GetContext(ctx, &current, `
		UPDATE application_statuses
		SET code = $1,
		    name = $2,
		    is_final = $3
		WHERE id = $4
		RETURNING id, code, name, is_final
	`, in.Code, in.Name, in.IsFinal, in.ID); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return dbmodel.ApplicationStatusDB{}, ErrStatusExists
		}
		return dbmodel.ApplicationStatusDB{}, err
	}

	if err = tx.Commit(); err != nil {
		return dbmodel.ApplicationStatusDB{}, err
	}
	return current, nil
}

func (r *AdminPostgres) DeleteStatus(ctx context.Context, statusID int16) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	var code string
	if err = tx.GetContext(ctx, &code, `SELECT code FROM application_statuses WHERE id=$1 FOR UPDATE`, statusID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrStatusNotFound
		}
		return err
	}

	if isProtectedStatusCode(code) {
		return ErrStatusProtected
	}

	if _, err = tx.ExecContext(ctx, `DELETE FROM application_statuses WHERE id=$1`, statusID); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23503" {
			return ErrStatusInUse
		}
		return err
	}

	return tx.Commit()
}

func isOperatorUserTx(ctx context.Context, tx *sqlx.Tx, userID int64) (bool, error) {
	var exists bool
	q := `
		SELECT EXISTS(
			SELECT 1
			FROM users u
			JOIN user_roles ur ON ur.user_id = u.id
			JOIN roles r ON r.id = ur.role_id
			WHERE u.id = $1 AND r.code = 'operator'
		)
	`
	if err := tx.GetContext(ctx, &exists, q, userID); err != nil {
		return false, err
	}
	return exists, nil
}

func currentStatusIDForAssign(ctx context.Context, tx *sqlx.Tx, code string) (int16, error) {
	return getStatusIDByCodeTx(ctx, tx, code)
}

func isProtectedStatusCode(code string) bool {
	switch code {
	case "new", "in_progress", "waiting", "resolved", "closed":
		return true
	default:
		return false
	}
}

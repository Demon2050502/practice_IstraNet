package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	dbmodel "practice_IstraNet/pkg/DB_model"
	"practice_IstraNet/pkg/dto"

	"github.com/jmoiron/sqlx"
)

type appState struct {
	AssignedTo *int64 `db:"assigned_to"`
	StatusCode string `db:"status_code"`
	IsFinal    bool   `db:"is_final"`
}

func (r *ApplicationPostgres) GetApplicationHistory(ctx context.Context, appID int64) ([]dbmodel.ApplicationHistoryDB, error) {
	var items []dbmodel.ApplicationHistoryDB

	q := `
		SELECT
			h.id,
			h.action,
			h.field,
			h.old_value,
			h.new_value,
			u.id AS actor_id,
			u.full_name AS actor_name,
			h.created_at
		FROM application_history h
		JOIN users u ON u.id = h.actor_id
		WHERE h.application_id = $1
		ORDER BY h.created_at ASC, h.id ASC
	`
	if err := r.db.SelectContext(ctx, &items, q, appID); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *ApplicationPostgres) TakeApplication(ctx context.Context, operatorID, appID int64) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	state, err := lockApplicationState(ctx, tx, appID)
	if err != nil {
		return err
	}
	if state.IsFinal {
		return ErrInvalidStatusTransition
	}
	if state.AssignedTo != nil {
		return ErrAlreadyAssigned
	}

	inProgressID, err := getStatusIDByCodeTx(ctx, tx, "in_progress")
	if err != nil {
		return err
	}

	if _, err = tx.ExecContext(ctx, `
		UPDATE applications
		SET assigned_to = $1,
		    status_id = $2,
		    updated_at = now()
		WHERE id = $3
	`, operatorID, inProgressID, appID); err != nil {
		return err
	}

	operatorIDStr := fmt.Sprintf("%d", operatorID)
	if err = insertHistoryTx(ctx, tx, appID, operatorID, "assign", strPtr("assigned_to"), nil, &operatorIDStr); err != nil {
		return err
	}

	if state.StatusCode != "in_progress" {
		if err = insertHistoryTx(ctx, tx, appID, operatorID, "status_change", strPtr("status_id"), strPtr(state.StatusCode), strPtr("in_progress")); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *ApplicationPostgres) ChangeApplicationStatus(ctx context.Context, operatorID int64, in dto.ChangeStatusRequest) error {
	if !isAllowedTargetStatus(in.StatusCode) {
		return ErrInvalidStatusCode
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

	state, err := lockApplicationState(ctx, tx, in.ID)
	if err != nil {
		return err
	}
	if state.AssignedTo == nil || *state.AssignedTo != operatorID {
		return ErrForbidden
	}
	if state.IsFinal {
		return ErrInvalidStatusTransition
	}
	if !isAllowedTransition(state.StatusCode, in.StatusCode) {
		return ErrInvalidStatusTransition
	}

	statusID, err := getStatusIDByCodeTx(ctx, tx, in.StatusCode)
	if err != nil {
		return err
	}

	if _, err = tx.ExecContext(ctx, `
		UPDATE applications
		SET status_id = $1,
		    updated_at = now()
		WHERE id = $2
	`, statusID, in.ID); err != nil {
		return err
	}

	if err = insertHistoryTx(ctx, tx, in.ID, operatorID, "status_change", strPtr("status_id"), strPtr(state.StatusCode), &in.StatusCode); err != nil {
		return err
	}

	if in.Comment != nil {
		if _, err = tx.ExecContext(ctx, `
			INSERT INTO application_comments(application_id, author_id, body)
			VALUES ($1,$2,$3)
		`, in.ID, operatorID, *in.Comment); err != nil {
			return err
		}
		if err = insertHistoryTx(ctx, tx, in.ID, operatorID, "comment", nil, nil, in.Comment); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *ApplicationPostgres) CloseApplication(ctx context.Context, operatorID int64, in dto.CloseApplicationRequest) error {
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
	if state.AssignedTo == nil || *state.AssignedTo != operatorID {
		return ErrForbidden
	}
	if state.StatusCode != "resolved" {
		return ErrInvalidStatusTransition
	}

	closedID, err := getStatusIDByCodeTx(ctx, tx, "closed")
	if err != nil {
		return err
	}

	if _, err = tx.ExecContext(ctx, `
		UPDATE applications
		SET status_id = $1,
		    updated_at = now(),
		    closed_at = now()
		WHERE id = $2
	`, closedID, in.ID); err != nil {
		return err
	}

	if err = insertHistoryTx(ctx, tx, in.ID, operatorID, "close", strPtr("status_id"), strPtr("resolved"), strPtr("closed")); err != nil {
		return err
	}

	if in.Comment != nil {
		if _, err = tx.ExecContext(ctx, `
			INSERT INTO application_comments(application_id, author_id, body)
			VALUES ($1,$2,$3)
		`, in.ID, operatorID, *in.Comment); err != nil {
			return err
		}
		if err = insertHistoryTx(ctx, tx, in.ID, operatorID, "comment", nil, nil, in.Comment); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func insertHistoryTx(ctx context.Context, tx *sqlx.Tx, appID, actorID int64, action string, field, oldValue, newValue *string) error {
	_, err := tx.ExecContext(ctx, `
		INSERT INTO application_history(application_id, actor_id, action, field, old_value, new_value)
		VALUES ($1,$2,$3,$4,$5,$6)
	`, appID, actorID, action, field, oldValue, newValue)
	return err
}

func lockApplicationState(ctx context.Context, tx *sqlx.Tx, appID int64) (appState, error) {
	var state appState
	err := tx.GetContext(ctx, &state, `
		SELECT a.assigned_to, s.code AS status_code, s.is_final
		FROM applications a
		JOIN application_statuses s ON s.id = a.status_id
		WHERE a.id = $1
		FOR UPDATE
	`, appID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return appState{}, ErrAppNotFound
		}
		return appState{}, err
	}
	return state, nil
}

func getStatusIDByCodeTx(ctx context.Context, tx *sqlx.Tx, code string) (int16, error) {
	var statusID int16
	if err := tx.GetContext(ctx, &statusID, `SELECT id FROM application_statuses WHERE code=$1`, code); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrStatusNotFound
		}
		return 0, err
	}
	return statusID, nil
}

func isAllowedTargetStatus(code string) bool {
	return code == "in_progress" || code == "waiting" || code == "resolved"
}

func isAllowedTransition(from, to string) bool {
	switch from {
	case "in_progress":
		return to == "waiting" || to == "resolved"
	case "waiting":
		return to == "in_progress" || to == "resolved"
	default:
		return false
	}
}


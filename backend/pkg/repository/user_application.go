package repository

import (
	"context"
	"database/sql"
	"errors"

	dbmodel "practice_IstraNet/pkg/DB_model"
	"practice_IstraNet/pkg/dto"
)

func (r *ApplicationPostgres) GetUserApplications(ctx context.Context, userID int64) ([]dbmodel.ApplicationDB, error) {
	var apps []dbmodel.ApplicationDB

	query := `
		SELECT a.id, a.title,
		       s.code AS status_code,
		       p.code AS priority_code,
		       a.created_at
		FROM applications a
		JOIN application_statuses s ON s.id = a.status_id
		JOIN application_priorities p ON p.id = a.priority_id
		WHERE a.created_by = $1
		ORDER BY a.created_at DESC
	`

	err := r.db.SelectContext(ctx, &apps, query, userID)
	return apps, err
}

func (r *ApplicationPostgres) GetUserApplicationByID(ctx context.Context, userID, appID int64) (dbmodel.ApplicationDB, error) {
	var app dbmodel.ApplicationDB

	query := `
		SELECT
			a.id, a.title, a.description,
			s.code AS status_code, s.name AS status_name,
			p.code AS priority_code, p.name AS priority_name, p.weight AS priority_weight,
			a.category_id,
			c.name AS category_name,
			u.id AS created_by_id,
			u.full_name AS created_by_name,
			a.assigned_to AS assigned_to_id,
			au.full_name AS assigned_to_name,
			a.contact_phone, a.contact_address,
			a.created_at, a.updated_at, a.closed_at
		FROM applications a
		JOIN application_statuses s ON s.id = a.status_id
		JOIN application_priorities p ON p.id = a.priority_id
		JOIN users u ON u.id = a.created_by
		LEFT JOIN application_categories c ON c.id = a.category_id
		LEFT JOIN users au ON au.id = a.assigned_to
		WHERE a.id = $1 AND a.created_by = $2
	`

	err := r.db.GetContext(ctx, &app, query, appID, userID)
	if err != nil {
		return dbmodel.ApplicationDB{}, ErrAppNotFound
	}
	return app, nil
}

func (r *ApplicationPostgres) DeleteUserApplication(ctx context.Context, userID, appID int64) error {
	res, err := r.db.ExecContext(ctx,
		`DELETE FROM applications 
		 WHERE id=$1 AND created_by=$2 
		   AND status_id IN (SELECT id FROM application_statuses WHERE is_final=FALSE)`,
		appID, userID)

	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return ErrForbidden
	}

	return nil
}

func (r *ApplicationPostgres) UpdateUserApplication(ctx context.Context, userID int64, in dto.ChangeApplicationRequest) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	var createdBy int64
	var oldTitle string
	var oldDescription string
	err = tx.QueryRowContext(ctx, `
		SELECT created_by, title, description
		FROM applications
		WHERE id = $1
		FOR UPDATE
	`, in.ID).Scan(&createdBy, &oldTitle, &oldDescription)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrAppNotFound
		}
		return err
	}
	if createdBy != userID {
		return ErrForbidden
	}

	newTitle := oldTitle
	newDescription := oldDescription
	if in.Title != nil {
		newTitle = *in.Title
	}
	if in.Description != nil {
		newDescription = *in.Description
	}

	if in.Title != nil || in.Description != nil {
		if _, err = tx.ExecContext(ctx, `
			UPDATE applications
			SET title = $1,
			    description = $2,
			    updated_at = now()
			WHERE id = $3
		`, newTitle, newDescription, in.ID); err != nil {
			return err
		}
	}

	if in.Title != nil && oldTitle != newTitle {
		if err = insertHistoryTx(ctx, tx, in.ID, userID, "edit", strPtr("title"), strPtr(oldTitle), strPtr(newTitle)); err != nil {
			return err
		}
	}
	if in.Description != nil && oldDescription != newDescription {
		if err = insertHistoryTx(ctx, tx, in.ID, userID, "edit", strPtr("description"), strPtr(oldDescription), strPtr(newDescription)); err != nil {
			return err
		}
	}

	if in.Comment != nil {
		if _, err = tx.ExecContext(ctx, `
			INSERT INTO application_comments(application_id, author_id, body)
			VALUES ($1,$2,$3)
		`, in.ID, userID, *in.Comment); err != nil {
			return err
		}

		if err = insertHistoryTx(ctx, tx, in.ID, userID, "comment", nil, nil, in.Comment); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *ApplicationPostgres) GetApplicationComments(ctx context.Context, appID int64) ([]dbmodel.ApplicationCommentDB, error) {
	var comments []dbmodel.ApplicationCommentDB

	query := `
		SELECT c.id, u.full_name AS author, c.body, c.created_at
		FROM application_comments c
		JOIN users u ON u.id = c.author_id
		WHERE c.application_id = $1
		ORDER BY c.created_at ASC
	`

	err := r.db.SelectContext(ctx, &comments, query, appID)
	return comments, err
}

func strPtr(v string) *string {
	return &v
}


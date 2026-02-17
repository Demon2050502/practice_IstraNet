package repository

import (
	"context"
	"errors"

	dbmodel "practice_IstraNet/pkg/DB_model"
	"practice_IstraNet/pkg/dto"
)

var ErrAppNotFound = errors.New("application not found")
var ErrForbidden = errors.New("forbidden")

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
	ORDER BY a.created_at DESC`

	err := r.db.SelectContext(ctx, &apps, query, userID)
	return apps, err
}

func (r *ApplicationPostgres) GetUserApplicationByID(ctx context.Context, userID, appID int64) (dbmodel.ApplicationDB, error) {
	var app dbmodel.ApplicationDB

	query := `
	SELECT a.*, 
	       s.code AS status_code, s.name AS status_name,
	       p.code AS priority_code, p.name AS priority_name, p.weight AS priority_weight,
	       u.full_name AS created_by_name
	FROM applications a
	JOIN application_statuses s ON s.id = a.status_id
	JOIN application_priorities p ON p.id = a.priority_id
	JOIN users u ON u.id = a.created_by
	WHERE a.id = $1 AND a.created_by = $2`

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

	// обновление title/description
	if in.Title != nil || in.Description != nil {
		_, err = tx.Exec(`
			UPDATE applications
			SET title = COALESCE($1, title),
			    description = COALESCE($2, description),
			    updated_at = now()
			WHERE id = $3 AND created_by = $4`,
			in.Title, in.Description, in.ID, userID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// добавление комментария
	if in.Comment != nil {
		_, err = tx.Exec(`
			INSERT INTO application_comments(application_id, author_id, body)
			VALUES ($1,$2,$3)`,
			in.ID, userID, *in.Comment)
		if err != nil {
			tx.Rollback()
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
	ORDER BY c.created_at ASC`

	err := r.db.SelectContext(ctx, &comments, query, appID)
	return comments, err
}

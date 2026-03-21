package repository

import (
	"context"
	"database/sql"
	"errors"

	dbmodel "practice_IstraNet/pkg/DB_model"

	"github.com/jmoiron/sqlx"
)

type ApplicationPostgres struct {
	db *sqlx.DB
}

func NewApplicationPostgres(db *sqlx.DB) *ApplicationPostgres {
	return &ApplicationPostgres{db: db}
}

func (r *ApplicationPostgres) CreateApplication(
	ctx context.Context,
	createdBy int64,
	title, description string,
	priorityCode string,
	categoryID *int64,
	contactPhone, contactAddress *string,
) (dbmodel.ApplicationDB, error) {

	if priorityCode == "" {
		priorityCode = "normal"
	}

	var statusID int16
	if err := r.db.GetContext(ctx, &statusID, `SELECT id FROM application_statuses WHERE code='new'`); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dbmodel.ApplicationDB{}, ErrStatusNotFound
		}
		return dbmodel.ApplicationDB{}, err
	}

	var priorityID int16
	if err := r.db.GetContext(ctx, &priorityID, `SELECT id FROM application_priorities WHERE code=$1`, priorityCode); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dbmodel.ApplicationDB{}, ErrPriorityNotFound
		}
		return dbmodel.ApplicationDB{}, err
	}

	if categoryID != nil {
		var exists bool
		if err := r.db.GetContext(ctx, &exists, `SELECT EXISTS(SELECT 1 FROM application_categories WHERE id=$1)`, *categoryID); err != nil {
			return dbmodel.ApplicationDB{}, err
		}
		if !exists {
			return dbmodel.ApplicationDB{}, ErrCategoryNotFound
		}
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return dbmodel.ApplicationDB{}, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	var appID int64
	qIns := `
		INSERT INTO applications (title, description, status_id, priority_id, category_id, created_by, contact_phone, contact_address)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		RETURNING id
	`
	if err = tx.QueryRowContext(ctx, qIns,
		title, description, statusID, priorityID, categoryID, createdBy, contactPhone, contactAddress,
	).Scan(&appID); err != nil {
		return dbmodel.ApplicationDB{}, err
	}

	if _, err = tx.ExecContext(ctx, `
		INSERT INTO application_history(application_id, actor_id, action)
		VALUES ($1, $2, 'create')
	`, appID, createdBy); err != nil {
		return dbmodel.ApplicationDB{}, err
	}

	if err = tx.Commit(); err != nil {
		return dbmodel.ApplicationDB{}, err
	}

	return r.getApplicationByID(ctx, appID)
}

func (r *ApplicationPostgres) GetAllApplications(ctx context.Context) ([]dbmodel.ApplicationDB, error) {
	var items []dbmodel.ApplicationDB

	q := `
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
		ORDER BY a.created_at DESC
	`

	if err := r.db.SelectContext(ctx, &items, q); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *ApplicationPostgres) getApplicationByID(ctx context.Context, id int64) (dbmodel.ApplicationDB, error) {
	var a dbmodel.ApplicationDB

	q := `
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
		WHERE a.id = $1
	`

	if err := r.db.GetContext(ctx, &a, q, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dbmodel.ApplicationDB{}, ErrAppNotFound
		}
		return dbmodel.ApplicationDB{}, err
	}

	return a, nil
}

func (r *ApplicationPostgres) GetApplicationByID(ctx context.Context, appID int64) (dbmodel.ApplicationDB, error) {
	return r.getApplicationByID(ctx, appID)
}

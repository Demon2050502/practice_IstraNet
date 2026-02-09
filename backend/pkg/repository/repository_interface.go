package repository

import (
	"context"

	"github.com/jmoiron/sqlx"

	dbmodel "practice_IstraNet/pkg/DB_model"
)

type Authorization interface{
	CreateUser(ctx context.Context, email, passwordHash, fullName, roleCode string) (userID int64, roleOut string, err error)
	GetUserByEmail(ctx context.Context, email string) (dbmodel.UserDB, error)
	GetUserRoleCode(ctx context.Context, userID int64) (string, error)
}

type Applications interface {
	CreateApplication(
		ctx context.Context,
		createdBy int64,
		title, description string,
		priorityCode string,
		categoryID *int64,
		contactPhone, contactAddress *string,
	) (dbmodel.ApplicationDB, error)

	GetAllApplications(ctx context.Context) ([]dbmodel.ApplicationDB, error)
}


type Repository struct {
	Authorization
	Applications
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Applications:  NewApplicationPostgres(db),
	}
}
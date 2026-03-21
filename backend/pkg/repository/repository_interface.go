package repository

import (
	"context"

	"github.com/jmoiron/sqlx"

	dbmodel "practice_IstraNet/pkg/DB_model"
	dto "practice_IstraNet/pkg/dto"
)

type Authorization interface {
	CreateUser(ctx context.Context, email, passwordHash, fullName, roleCode string) (userID int64, roleOut string, err error)
	GetUserByEmail(ctx context.Context, email string) (dbmodel.UserDB, error)
	GetUserRoleCode(ctx context.Context, userID int64) (string, error)
}

type Administration interface {
	AssignApplication(ctx context.Context, adminID int64, in dto.AdminAssignApplicationRequest) error
	ChangeApplicationStatusByAdmin(ctx context.Context, adminID int64, in dto.AdminChangeApplicationStatusRequest) error
	DeleteApplicationByAdmin(ctx context.Context, appID int64) error
	GetUsers(ctx context.Context) ([]dbmodel.AdminUserDB, error)
	GetUserByIDForAdmin(ctx context.Context, userID int64) (dbmodel.AdminUserDB, error)
	ChangeUserRole(ctx context.Context, adminID int64, in dto.AdminChangeUserRoleRequest) error
	DeleteUserByAdmin(ctx context.Context, adminID, userID int64) error
	CreateStatus(ctx context.Context, in dto.AdminCreateStatusRequest) (dbmodel.ApplicationStatusDB, error)
	UpdateStatus(ctx context.Context, in dto.AdminUpdateStatusRequest) (dbmodel.ApplicationStatusDB, error)
	DeleteStatus(ctx context.Context, statusID int16) error
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
	GetApplicationByID(ctx context.Context, appID int64) (dbmodel.ApplicationDB, error)

	DeleteUserApplication(ctx context.Context, userID, appID int64) error
	GetUserApplications(ctx context.Context, userID int64) ([]dbmodel.ApplicationDB, error)
	GetUserApplicationByID(ctx context.Context, userID, appID int64) (dbmodel.ApplicationDB, error)
	UpdateUserApplication(ctx context.Context, userID int64, in dto.ChangeApplicationRequest) error
	GetApplicationComments(ctx context.Context, appID int64) ([]dbmodel.ApplicationCommentDB, error)
	GetApplicationHistory(ctx context.Context, appID int64) ([]dbmodel.ApplicationHistoryDB, error)
	TakeApplication(ctx context.Context, operatorID, appID int64) error
	ChangeApplicationStatus(ctx context.Context, operatorID int64, in dto.ChangeStatusRequest) error
	CloseApplication(ctx context.Context, operatorID int64, in dto.CloseApplicationRequest) error
}

type Repository struct {
	Authorization
	Applications
	Administration
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization:  NewAuthPostgres(db),
		Applications:   NewApplicationPostgres(db),
		Administration: NewAdminPostgres(db),
	}
}

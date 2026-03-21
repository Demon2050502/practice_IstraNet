package service

import (
	"context"

	dto "practice_IstraNet/pkg/dto"
	repository "practice_IstraNet/pkg/repository"
)

type Authorization interface {
	SignUp(ctx context.Context, in dto.SignUpRequest) (dto.AuthResponse, error)
	SignIn(ctx context.Context, in dto.SignInRequest) (dto.AuthResponse, error)
	GenerateJWT(userID int64, name, role string) (string, error)
}

type Administration interface {
	AssignApplication(ctx context.Context, adminID int64, in dto.AdminAssignApplicationRequest) error
	ChangeApplicationStatusByAdmin(ctx context.Context, adminID int64, in dto.AdminChangeApplicationStatusRequest) error
	DeleteApplicationByAdmin(ctx context.Context, appID int64) error
	GetUsers(ctx context.Context) (dto.AdminUsersResponse, error)
	GetUserByIDForAdmin(ctx context.Context, userID int64) (dto.AdminUserResponse, error)
	ChangeUserRole(ctx context.Context, adminID int64, in dto.AdminChangeUserRoleRequest) error
	DeleteUserByAdmin(ctx context.Context, adminID, userID int64) error
	CreateStatus(ctx context.Context, in dto.AdminCreateStatusRequest) (dto.ApplicationStatusResponse, error)
	UpdateStatus(ctx context.Context, in dto.AdminUpdateStatusRequest) (dto.ApplicationStatusResponse, error)
	DeleteStatus(ctx context.Context, statusID int16) error
}

type Applications interface {
	CreateApplication(ctx context.Context, userID int64, in dto.CreateApplicationRequest) (dto.ApplicationResponse, error)
	GetAllApplications(ctx context.Context) (dto.ApplicationsListResponse, error)

	GetUserApplications(ctx context.Context, userID int64) (dto.UserAppsResponse, error)
	GetUserApplicationByID(ctx context.Context, userID, appID int64) (dto.UserAppDetailsResponse, error)
	DeleteUserApplication(ctx context.Context, userID, appID int64) error
	UpdateUserApplication(ctx context.Context, userID int64, in dto.ChangeApplicationRequest) error
	GetOperatorApplicationByID(ctx context.Context, appID int64) (dto.UserAppDetailsResponse, error)
	GetApplicationHistory(ctx context.Context, appID int64) (dto.ApplicationHistoryResponse, error)
	TakeApplication(ctx context.Context, operatorID, appID int64) error
	ChangeApplicationStatus(ctx context.Context, operatorID int64, in dto.ChangeStatusRequest) error
	CloseApplication(ctx context.Context, operatorID int64, in dto.CloseApplicationRequest) error
}

type Service struct {
	Authorization
	Applications
	Administration
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization:  NewAuthService(repos.Authorization),
		Applications:   NewApplicationService(repos.Applications),
		Administration: NewAdminService(repos.Administration),
	}
}

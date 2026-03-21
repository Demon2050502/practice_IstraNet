package service

import (
	"context"

	dbmodel "practice_IstraNet/pkg/DB_model"
	"practice_IstraNet/pkg/dto"
	repository "practice_IstraNet/pkg/repository"
)

type AdminService struct {
	repos repository.Administration
}

func NewAdminService(repos repository.Administration) *AdminService {
	return &AdminService{repos: repos}
}

func (s *AdminService) AssignApplication(ctx context.Context, adminID int64, in dto.AdminAssignApplicationRequest) error {
	return s.repos.AssignApplication(ctx, adminID, in)
}

func (s *AdminService) ChangeApplicationStatusByAdmin(ctx context.Context, adminID int64, in dto.AdminChangeApplicationStatusRequest) error {
	return s.repos.ChangeApplicationStatusByAdmin(ctx, adminID, in)
}

func (s *AdminService) DeleteApplicationByAdmin(ctx context.Context, appID int64) error {
	return s.repos.DeleteApplicationByAdmin(ctx, appID)
}

func (s *AdminService) GetUsers(ctx context.Context) (dto.AdminUsersResponse, error) {
	items, err := s.repos.GetUsers(ctx)
	if err != nil {
		return dto.AdminUsersResponse{}, err
	}

	out := dto.AdminUsersResponse{
		Items: make([]dto.AdminUserResponse, 0, len(items)),
	}
	for _, item := range items {
		out.Items = append(out.Items, mapAdminUserDB(item))
	}
	return out, nil
}

func (s *AdminService) GetUserByIDForAdmin(ctx context.Context, userID int64) (dto.AdminUserResponse, error) {
	item, err := s.repos.GetUserByIDForAdmin(ctx, userID)
	if err != nil {
		return dto.AdminUserResponse{}, err
	}
	return mapAdminUserDB(item), nil
}

func (s *AdminService) ChangeUserRole(ctx context.Context, adminID int64, in dto.AdminChangeUserRoleRequest) error {
	return s.repos.ChangeUserRole(ctx, adminID, in)
}

func (s *AdminService) DeleteUserByAdmin(ctx context.Context, adminID, userID int64) error {
	return s.repos.DeleteUserByAdmin(ctx, adminID, userID)
}

func (s *AdminService) CreateStatus(ctx context.Context, in dto.AdminCreateStatusRequest) (dto.ApplicationStatusResponse, error) {
	item, err := s.repos.CreateStatus(ctx, in)
	if err != nil {
		return dto.ApplicationStatusResponse{}, err
	}
	return mapStatusDB(item), nil
}

func (s *AdminService) UpdateStatus(ctx context.Context, in dto.AdminUpdateStatusRequest) (dto.ApplicationStatusResponse, error) {
	item, err := s.repos.UpdateStatus(ctx, in)
	if err != nil {
		return dto.ApplicationStatusResponse{}, err
	}
	return mapStatusDB(item), nil
}

func (s *AdminService) DeleteStatus(ctx context.Context, statusID int16) error {
	return s.repos.DeleteStatus(ctx, statusID)
}

func mapAdminUserDB(in dbmodel.AdminUserDB) dto.AdminUserResponse {
	return dto.AdminUserResponse{
		ID:        in.ID,
		Email:     in.Email,
		FullName:  in.FullName,
		RoleCode:  in.RoleCode,
		RoleName:  in.RoleName,
		IsActive:  in.IsActive,
		CreatedAt: in.CreatedAt,
	}
}

func mapStatusDB(in dbmodel.ApplicationStatusDB) dto.ApplicationStatusResponse {
	return dto.ApplicationStatusResponse{
		ID:      in.ID,
		Code:    in.Code,
		Name:    in.Name,
		IsFinal: in.IsFinal,
	}
}

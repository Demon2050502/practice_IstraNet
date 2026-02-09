package service

import (
	"context"

	dto "practice_IstraNet/pkg/dto"
	repository "practice_IstraNet/pkg/repository"
)

type Authorization interface{
	SignUp(ctx context.Context, in dto.SignUpRequest) (dto.AuthResponse, error)
	SignIn(ctx context.Context, in dto.SignInRequest) (dto.AuthResponse, error)
	GenerateJWT(userID int64, name, role string) (string, error)
}

type Applications interface {
	CreateApplication(ctx context.Context, userID int64, in dto.CreateApplicationRequest) (dto.ApplicationResponse, error)
	GetAllApplications(ctx context.Context) (dto.ApplicationsListResponse, error)
}

type Service struct {
	Authorization
	Applications
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Applications:  NewApplicationService(repos.Applications),
	}
}
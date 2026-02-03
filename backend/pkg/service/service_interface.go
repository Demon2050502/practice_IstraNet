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

type Service struct {
	Authorization
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
	}
}
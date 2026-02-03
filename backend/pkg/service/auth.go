package service

import (
	"context"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	dto "practice_IstraNet/pkg/dto"
	repository "practice_IstraNet/pkg/repository"
)

type AuthService struct {
	repos repository.Authorization 
}

func NewAuthService(repos repository.Authorization) *AuthService {
	return &AuthService{repos: repos}
}


func (s *AuthService) SignUp(ctx context.Context, in dto.SignUpRequest) (dto.AuthResponse, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	userID, role, err := s.repos.CreateUser(ctx, in.Email, string(hash), in.FullName, in.Role)
	if err != nil {
		return dto.AuthResponse{}, err
	} 

	token, err := s.GenerateJWT(userID, in.FullName, role)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	return dto.AuthResponse{
		Token: token,
		User: dto.AuthUser{
			ID:   userID,
			Name: in.FullName,
			Role: role,
		},
	}, nil
}

func (s *AuthService) SignIn(ctx context.Context, in dto.SignInRequest) (dto.AuthResponse, error) {
	u, err := s.repos.GetUserByEmail(ctx, in.Email)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(in.Password)); err != nil {
		return dto.AuthResponse{}, dto.ErrInvalidCredentials
	}

	role, err := s.repos.GetUserRoleCode(ctx, u.ID)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	token, err := s.GenerateJWT(u.ID, u.FullName, role)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	return dto.AuthResponse{
		Token: token,
		User: dto.AuthUser{
			ID:   u.ID,
			Name: u.FullName,
			Role: role,
		},
	}, nil
}

func (s *AuthService) GenerateJWT(userID int64, name, role string) (string, error) {
	now := time.Now()
	jwtSecret := getSecretKey()

	claims := jwt.MapClaims{
		"sub":  userID,
		"name": name,
		"role": role,
		"iat":  now.Unix(),
		"exp":  now.Add(24 * time.Hour).Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return t.SignedString(jwtSecret)
}

func getSecretKey() ([]byte) {
	jwtSecret := os.Getenv("JWT_SECRET")
	return []byte(jwtSecret)
}
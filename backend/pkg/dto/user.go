package dto

import "errors"

type SignUpRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
    FullName string `json:"full_name" binding:"required"`
    Role string `json:"role"`
}

type SignInRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

type AuthUser struct {
    ID   int64  `json:"id"`
    Name string `json:"name"`
    Role string `json:"role"`
}

type AuthResponse struct {
    Token string   `json:"token"`
    User  AuthUser `json:"user"`
}

type ErrorResponse struct {
    Error   string `json:"error"`
    Message string `json:"message,omitempty"`
}

var ErrInvalidCredentials = errors.New("invalid credentials")

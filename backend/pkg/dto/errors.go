package dto

import "errors"

var (
	// Ошибки аутентификации.
	ErrInvalidCredentials = errors.New("invalid credentials")
)

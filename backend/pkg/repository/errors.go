package repository

import "errors"

var (
	// Ошибки поиска заявок и связанных справочников.
	ErrAppNotFound      = errors.New("application not found")
	ErrPriorityNotFound = errors.New("priority not found")
	ErrCategoryNotFound = errors.New("category not found")
	ErrStatusNotFound   = errors.New("status not found")
	ErrOperatorNotFound = errors.New("operator not found")

	// Ошибки прав доступа и бизнес-правил по заявкам.
	ErrForbidden               = errors.New("forbidden")
	ErrAlreadyAssigned         = errors.New("application already assigned")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
	ErrInvalidStatusCode       = errors.New("invalid status code")
	ErrFinalApplication        = errors.New("final application")

	// Ошибки управления пользователями и ролями.
	ErrUserExists          = errors.New("user already exists")
	ErrUserNotFound        = errors.New("user not found")
	ErrUserInactive        = errors.New("user inactive")
	ErrRoleNotFound        = errors.New("role not found")
	ErrNoUserRole          = errors.New("user role not found")
	ErrSelfActionForbidden = errors.New("self action forbidden")
	ErrUserHasRelations    = errors.New("user has relations")

	// Ошибки управления справочниками.
	ErrStatusProtected = errors.New("status is protected")
	ErrStatusInUse     = errors.New("status is in use")
	ErrStatusExists    = errors.New("status already exists")
)

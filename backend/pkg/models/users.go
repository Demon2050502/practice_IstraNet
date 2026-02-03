package models

type UserModel struct {
	ID           int64
	Email        string
	PasswordHash string
	FullName     string
	IsActive     bool
	RoleID       int8
	RoleName     string
} // пока не использовал
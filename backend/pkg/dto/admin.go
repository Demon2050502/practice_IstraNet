package dto

import "time"

type AdminAssignApplicationRequest struct {
	ID         int64 `json:"id" binding:"required"`
	OperatorID int64 `json:"operator_id" binding:"required"`
}

type AdminChangeApplicationStatusRequest struct {
	ID         int64   `json:"id" binding:"required"`
	StatusCode string  `json:"status_code" binding:"required"`
	Comment    *string `json:"comment,omitempty"`
}

type AdminDeleteApplicationRequest struct {
	ID int64 `json:"id" binding:"required"`
}

type AdminUserResponse struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	RoleCode  string    `json:"role_code"`
	RoleName  string    `json:"role_name"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type AdminUsersResponse struct {
	Items []AdminUserResponse `json:"items"`
}

type AdminChangeUserRoleRequest struct {
	UserID   int64  `json:"user_id" binding:"required"`
	RoleCode string `json:"role_code" binding:"required"`
}

type AdminDeleteUserRequest struct {
	UserID int64 `json:"user_id" binding:"required"`
}

type ApplicationStatusResponse struct {
	ID      int16  `json:"id"`
	Code    string `json:"code"`
	Name    string `json:"name"`
	IsFinal bool   `json:"is_final"`
}

type AdminCreateStatusRequest struct {
	Code    string `json:"code" binding:"required"`
	Name    string `json:"name" binding:"required"`
	IsFinal bool   `json:"is_final"`
}

type AdminUpdateStatusRequest struct {
	ID      int16  `json:"id" binding:"required"`
	Code    string `json:"code" binding:"required"`
	Name    string `json:"name" binding:"required"`
	IsFinal bool   `json:"is_final"`
}

type AdminDeleteStatusRequest struct {
	ID int16 `json:"id" binding:"required"`
}

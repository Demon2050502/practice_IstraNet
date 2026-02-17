package dto

import "time"

type UserAppShort struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Status    string    `json:"status"`
	Priority  string    `json:"priority"`
	CreatedAt time.Time `json:"created_at"`
}

type UserAppsResponse struct {
	Items []UserAppShort `json:"items"`
}

type UserAppDetailsResponse struct {
	ApplicationResponse
	Comments []AppComment `json:"comments"`
}

type AppComment struct {
	ID        int64     `json:"id"`
	Author    string    `json:"author"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

type ChangeApplicationRequest struct {
	ID          int64   `json:"id" binding:"required"`
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Comment     *string `json:"comment,omitempty"`
}

type DeleteApplicationRequest struct {
	ID int64 `json:"id" binding:"required"`
}

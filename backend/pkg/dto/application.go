package dto

import "time"

type CreateApplicationRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`

	// Можно не передавать — тогда будет "normal"
	PriorityCode string `json:"priority_code,omitempty"` // low|normal|high|critical

	// category_id опционально
	CategoryID *int64 `json:"category_id,omitempty"`

	ContactPhone   *string `json:"contact_phone,omitempty"`
	ContactAddress *string `json:"contact_address,omitempty"`
}

type ApplicationResponse struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`

	Status struct {
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"status"`

	Priority struct {
		Code   string `json:"code"`
		Name   string `json:"name"`
		Weight int16  `json:"weight"`
	} `json:"priority"`

	Category *struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"category,omitempty"`

	CreatedBy struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"created_by"`

	AssignedTo *struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"assigned_to,omitempty"`

	ContactPhone   *string `json:"contact_phone,omitempty"`
	ContactAddress *string `json:"contact_address,omitempty"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	ClosedAt  *time.Time `json:"closed_at,omitempty"`
}

type ApplicationsListResponse struct {
	Items []ApplicationResponse `json:"items"`
}

package db_model

import "time"

type ApplicationDB struct {
	ID          int64  `db:"id"`
	Title       string `db:"title"`
	Description string `db:"description"`

	StatusCode string `db:"status_code"`
	StatusName string `db:"status_name"`

	PriorityCode   string `db:"priority_code"`
	PriorityName   string `db:"priority_name"`
	PriorityWeight int16  `db:"priority_weight"`

	CategoryID   *int64  `db:"category_id"`
	CategoryName *string `db:"category_name"`

	CreatedByID   int64  `db:"created_by_id"`
	CreatedByName string `db:"created_by_name"`

	AssignedToID   *int64  `db:"assigned_to_id"`
	AssignedToName *string `db:"assigned_to_name"`

	ContactPhone   *string `db:"contact_phone"`
	ContactAddress *string `db:"contact_address"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	ClosedAt  *time.Time `db:"closed_at"`
}

package db_model

import "time"

type AdminUserDB struct {
	ID        int64     `db:"id"`
	Email     string    `db:"email"`
	FullName  string    `db:"full_name"`
	RoleCode  string    `db:"role_code"`
	RoleName  string    `db:"role_name"`
	IsActive  bool      `db:"is_active"`
	CreatedAt time.Time `db:"created_at"`
}

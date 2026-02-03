package db_model

type UserDB struct {
	ID           int64  `db:"id"`
	Email        string `db:"email"`
	PasswordHash string `db:"password_hash"`
	FullName     string `db:"full_name"`
	IsActive     bool   `db:"is_active"`
}
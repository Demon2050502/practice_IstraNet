package db_model

import "time"

type ApplicationCommentDB struct {
	ID        int64     `db:"id"`
	Author    string    `db:"author"`
	Body      string    `db:"body"`
	CreatedAt time.Time `db:"created_at"`
}

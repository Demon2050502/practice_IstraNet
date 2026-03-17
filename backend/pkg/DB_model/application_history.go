package db_model

import "time"

type ApplicationHistoryDB struct {
	ID        int64     `db:"id"`
	Action    string    `db:"action"`
	Field     *string   `db:"field"`
	OldValue  *string   `db:"old_value"`
	NewValue  *string   `db:"new_value"`
	ActorID   int64     `db:"actor_id"`
	ActorName string    `db:"actor_name"`
	CreatedAt time.Time `db:"created_at"`
}


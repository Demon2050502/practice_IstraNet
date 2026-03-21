package db_model

type ApplicationStatusDB struct {
	ID      int16  `db:"id"`
	Code    string `db:"code"`
	Name    string `db:"name"`
	IsFinal bool   `db:"is_final"`
}

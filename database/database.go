package database

import "database/sql"

func Open() (DB *sql.DB, err error) {
	DB, err = sql.Open("pgx", "host=localhost port=5432 user=money-api password=money-api dbname=money-api sslmode=disable")
	if err != nil {
		return
	}
	_, err = DB.Exec("DELETE FROM users;")
	if err != nil {
		return
	}
	return
}

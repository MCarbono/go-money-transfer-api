package database

import "database/sql"

func Open() (DB *sql.DB, err error) {
	DB, err = sql.Open("pgx", "host=localhost port=5432 user=money-api password=root dbname=money-api sslmode=disable")
	if err != nil {
		return
	}
	err = DB.Ping()
	if err != nil {
		return
	}
	return
}

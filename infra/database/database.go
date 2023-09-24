package database

import (
	"database/sql"
	"fmt"
	"time"
)

const maxRetries = 5
const retryInterval = 2 * time.Second

func Open(host string) (DB *sql.DB, err error) {
	DB, err = sql.Open("pgx", fmt.Sprintf("host=%v port=5432 user=money-api password=root dbname=money-api sslmode=disable", host))
	if err != nil {
		return
	}
	for i := 0; i < maxRetries; i++ {
		err = DB.Ping()
		if err == nil {
			break
		}
		if err != nil {
			fmt.Printf("Connection failed (Attempt %d): %v\n", i+1, err)
			time.Sleep(retryInterval)
		}
	}
	fmt.Println("Connected to the database!")
	return
}

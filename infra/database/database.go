package database

import (
	"database/sql"
	"fmt"
	"time"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

const maxRetries = 5
const retryInterval = 2 * time.Second

func Open(dbConfig DatabaseConfig) (DB *sql.DB, err error) {
	DB, err = sql.Open(
		"pgx", fmt.Sprintf(
			"host=%v port=%v user=%v password=%v dbname=%v sslmode=disable",
			dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.Name),
	)
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

package main

import (
	"money-transfer-api/database"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	DB, err := database.Open()
	if err != nil {
		panic(err)
	}
	defer DB.Close()
}

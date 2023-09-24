package main

import (
	"fmt"
	"log"
	"money-transfer-api/infra/controller"
	"money-transfer-api/infra/database"
	"money-transfer-api/infra/router"
	"money-transfer-api/repository"
	"money-transfer-api/service"
	"money-transfer-api/uow"
	"net/http"

	_ "github.com/jackc/pgx/v4/stdlib"
)

const (
	PORT = 3000
)

func main() {
	DB, err := database.Open("db")
	if err != nil {
		panic(err)
	}
	defer DB.Close()
	uow := uow.NewUowImpl(DB)
	uow.Register("UserRepository", func() interface{} {
		repo := repository.NewUserRepositoryPostgres(uow.Db, uow.Tx)
		return repo
	})
	userService := service.NewUser(DB, uow)
	controller := controller.NewUserController(*userService)
	router := router.NewRouter(*controller)
	fmt.Printf("Starting server on port %v\n", PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", PORT), router))
}

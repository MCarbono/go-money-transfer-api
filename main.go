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
	DB, err := database.Open()
	if err != nil {
		panic(err)
	}
	defer DB.Close()
	uow := uow.NewUowImpl(DB)
	userService := service.NewUser(DB, uow, repository.USER_REPOSITORY_POSTGRES)
	controller := controller.NewUserController(*userService)
	router := router.NewRouter(*controller)
	fmt.Printf("Starting server on port %v\n", PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", PORT), router))
}

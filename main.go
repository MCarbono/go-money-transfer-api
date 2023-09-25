package main

import (
	"fmt"
	"log"
	"money-transfer-api/config"
	"money-transfer-api/infra/controller"
	"money-transfer-api/infra/database"
	"money-transfer-api/infra/router"
	"money-transfer-api/repository"
	"money-transfer-api/service"
	"money-transfer-api/uow"
	"net/http"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	config, err := config.LoadEnvConfig()
	if err != nil {
		panic(err)
	}
	DB, err := database.Open(config.DBHost)
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
	fmt.Printf("Starting server on port %v\n", config.ServerPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", config.ServerPort), router))
}

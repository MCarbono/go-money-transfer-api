package router

import (
	"money-transfer-api/infra/controller"
	"net/http"

	"github.com/go-chi/chi"
)

func NewRouter(userController controller.UserController) http.Handler {
	r := chi.NewRouter()
	r.Post("/transfer", userController.Transfer)
	r.Get("/balance/{id}", userController.GetBalance)
	return r
}

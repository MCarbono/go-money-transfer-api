package controller

import (
	"encoding/json"
	"fmt"
	"money-transfer-api/service"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type UserController struct {
	UserService service.User
}

func NewUserController(userService service.User) *UserController {
	return &UserController{
		UserService: userService,
	}
}

func (u *UserController) Transfer(w http.ResponseWriter, r *http.Request) {
	var input service.TransferInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrMessage{Err: "invalid payload"})
		return
	}
	err = u.UserService.Transfer(&input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrMessage{Err: err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (u *UserController) GetBalance(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrMessage{Err: "invalid userIDs"})
		return
	}
	balance, err := u.UserService.GetBalance(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrMessage{Err: "internal server error"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(balance)
}

type ErrMessage struct {
	Err string `json:"err"`
}

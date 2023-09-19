package test

import (
	"money-transfer-api/database"
	"money-transfer-api/service"
	"testing"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func TestUserService(t *testing.T) {
	DB, err := database.Open()
	if err != nil {
		panic(err)
	}
	defer DB.Close()

	t.Run("Should get the user balance by its ID", func(t *testing.T) {
		defer DB.Exec("DELETE FROM users;")
		_, err = DB.Exec("INSERT INTO users (id, username, balance) VALUES (1, 'first_user', 2000);")
		if err != nil {
			t.Error(err)
			return
		}
		user := service.NewUser(DB)
		output, err := user.GetBalance(1)
		if err != nil {
			t.Error(err)
			return
		}
		if output.Balance != 2000 {
			t.Errorf("User balance should be equal to 2000, but got %v", output.Balance)
		}
	})
}

// _, err = DB.Exec("INSERT INTO users (id, username, balance) VALUES (1, 'first_user', 2000);")
// 	if err != nil {
// 		return
// 	}
// 	_, err = DB.Exec("INSERT INTO users (id, username, balance) VALUES (2, 'second_user', 0);")
// 	if err != nil {
// 		return
// 	}
// 	_, err = DB.Exec("INSERT INTO users (id, username, balance) VALUES (3, 'third_user', 1000);")
// 	if err != nil {
// 		return
// 	}

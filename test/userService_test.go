package test

import (
	"context"
	"errors"
	"money-transfer-api/infra/database"
	"money-transfer-api/repository"
	"money-transfer-api/service"
	"money-transfer-api/test/dockertest"
	"money-transfer-api/uow"
	"sync"
	"testing"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func TestUserService(t *testing.T) {
	ctx := context.Background()
	container, err := dockertest.StartPostgresContainer(ctx)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := container.Terminate(ctx); err != nil {
			panic(err)
		}
	}()
	dbConfig := database.DatabaseConfig{
		Host:     container.Host,
		Port:     container.Port,
		User:     dockertest.DbUser,
		Password: dockertest.DbPassword,
		Name:     dockertest.DbName,
	}
	DB, err := database.Open(dbConfig)
	if err != nil {
		panic(err)
	}
	defer DB.Close()
	_, err = DB.Exec("create table users(id integer primary key, username text, balance numeric);")
	if err != nil {
		panic(err)
	}
	uowWithUserRepo := uow.NewUowImpl(DB, repository.NewUserRepositoryPostgres())

	t.Run("Should get the user balance by its ID", func(t *testing.T) {
		DB.Exec("DELETE FROM users;")
		defer DB.Exec("DELETE FROM users;")
		_, err = DB.Exec("INSERT INTO users (id, username, balance) VALUES (1, 'first_user', 2000);")
		if err != nil {
			t.Error(err)
			return
		}
		user := service.NewUser(DB, uowWithUserRepo)
		output, err := user.GetBalance(1)
		if err != nil {
			t.Error(err)
			return
		}
		want := 2000.0
		if output.Balance != want {
			t.Errorf("User balance should be equal to %v, but got %v", want, output.Balance)
		}
	})

	t.Run("Should be able make 5 transfers on sequence from one user to another", func(t *testing.T) {
		DB.Exec("DELETE FROM users;")
		defer DB.Exec("DELETE FROM users;")
		_, err = DB.Exec("INSERT INTO users (id, username, balance) VALUES (1, 'first_user', 500);")
		if err != nil {
			t.Error(err)
			return
		}
		_, err = DB.Exec("INSERT INTO users (id, username, balance) VALUES (2, 'second_user', 0);")
		if err != nil {
			t.Error(err)
			return
		}
		user := service.NewUser(DB, uowWithUserRepo)
		var wg sync.WaitGroup
		totalRequests := 5
		wg.Add(totalRequests)
		for i := 0; i < totalRequests; i++ {
			go func() {
				defer wg.Done()
				err = user.Transfer(&service.TransferInput{
					Amount:        100,
					DebtorID:      1,
					BeneficiaryID: 2,
				})
				if err != nil {
					t.Errorf("%T", err)
					t.Error(err)
				}
			}()
		}
		wg.Wait()
		outputDebtor, err := user.GetBalance(1)
		if err != nil {
			t.Error(err)
			return
		}
		outputBeneficiary, err := user.GetBalance(2)
		if err != nil {
			t.Error(err)
			return
		}
		outputDebtorWant := 0.0
		if outputDebtor.Balance != outputDebtorWant {
			t.Errorf("Transfer() failed. Debtor balance want: %v, got %v", outputDebtorWant, outputDebtor.Balance)
		}
		outputBeneficiaryWant := 500.0
		if outputBeneficiary.Balance != outputBeneficiaryWant {
			t.Errorf("Transfer() failed. Beneciary balance want: %v, got %v", outputBeneficiaryWant, outputBeneficiary.Balance)
		}
	})

	t.Run("Should be able make 2 transfers from differents users to a third one", func(t *testing.T) {
		DB.Exec("DELETE FROM users;")
		defer DB.Exec("DELETE FROM users;")
		_, err = DB.Exec("INSERT INTO users (id, username, balance) VALUES (1, 'first_user', 500);")
		if err != nil {
			t.Error(err)
			return
		}
		_, err = DB.Exec("INSERT INTO users (id, username, balance) VALUES (2, 'second_user', 0);")
		if err != nil {
			t.Error(err)
			return
		}
		_, err = DB.Exec("INSERT INTO users (id, username, balance) VALUES (3, 'third_user', 1000);")
		if err != nil {
			t.Error(err)
			return
		}
		user := service.NewUser(DB, uowWithUserRepo)
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			err := user.Transfer(&service.TransferInput{
				Amount:        100,
				DebtorID:      1,
				BeneficiaryID: 2,
			})
			if err != nil {
				t.Errorf("%T", err)
				t.Error(err)
			}
		}()
		go func() {
			defer wg.Done()
			err := user.Transfer(&service.TransferInput{
				Amount:        100,
				DebtorID:      3,
				BeneficiaryID: 2,
			})
			if err != nil {
				t.Errorf("%T", err)
				t.Error(err)
			}
		}()
		wg.Wait()
		outputFirstDebtor, err := user.GetBalance(1)
		if err != nil {
			t.Error(err)
			return
		}
		outputSecondDebtor, err := user.GetBalance(3)
		if err != nil {
			t.Error(err)
			return
		}
		outputBeneficiary, err := user.GetBalance(2)
		if err != nil {
			t.Error(err)
			return
		}
		outputDebtorFirstWant := 400.0
		if outputFirstDebtor.Balance != outputDebtorFirstWant {
			t.Errorf("Transfer() failed. want: %v, got %v", outputDebtorFirstWant, outputFirstDebtor.Balance)
		}
		outputDebtorSecondWant := 900.0
		if outputSecondDebtor.Balance != outputDebtorSecondWant {
			t.Errorf("Transfer() failed. want: %v, got %v", outputDebtorSecondWant, outputSecondDebtor.Balance)
		}
		outputBeneficiaryWant := 200.0
		if outputBeneficiary.Balance != outputBeneficiaryWant {
			t.Errorf("Transfer() failed. want: %v, got %v", outputBeneficiaryWant, outputBeneficiary.Balance)
		}
	})

	t.Run("Should not be able to make a transfer because debtor user does not have enough balance", func(t *testing.T) {
		DB.Exec("DELETE FROM users;")
		defer DB.Exec("DELETE FROM users;")
		_, err = DB.Exec("INSERT INTO users (id, username, balance) VALUES (1, 'first_user', 500);")
		if err != nil {
			t.Error(err)
			return
		}
		_, err = DB.Exec("INSERT INTO users (id, username, balance) VALUES (2, 'second_user', 0);")
		if err != nil {
			t.Error(err)
			return
		}
		user := service.NewUser(DB, uowWithUserRepo)
		got := user.Transfer(&service.TransferInput{
			Amount:        1000,
			DebtorID:      1,
			BeneficiaryID: 2,
		})
		want := errors.New("insufficient funds")
		if want.Error() != got.Error() {
			t.Errorf("Transfer() failed. want %v, got %v", want, got)
		}
	})

	t.Run("Should be able make 5 transfers to user 1 to 2 and 5 transfers to user 2 to 1", func(t *testing.T) {
		DB.Exec("DELETE FROM users;")
		defer DB.Exec("DELETE FROM users;")
		_, err = DB.Exec("INSERT INTO users (id, username, balance) VALUES (1, 'first_user', 500);")
		if err != nil {
			t.Error(err)
			return
		}
		_, err = DB.Exec("INSERT INTO users (id, username, balance) VALUES (2, 'second_user', 500);")
		if err != nil {
			t.Error(err)
			return
		}
		user := service.NewUser(DB, uowWithUserRepo)
		var wg sync.WaitGroup
		totalRequests := 10
		wg.Add(totalRequests)
		for i := 0; i < totalRequests; i++ {
			if i%2 == 0 {
				go func() {
					defer wg.Done()
					err = user.Transfer(&service.TransferInput{
						Amount:        100,
						DebtorID:      1,
						BeneficiaryID: 2,
					})
					if err != nil {
						t.Log(err)
					}
				}()
			} else {
				go func() {
					defer wg.Done()
					err = user.Transfer(&service.TransferInput{
						Amount:        50,
						DebtorID:      2,
						BeneficiaryID: 1,
					})
					if err != nil {
						t.Log(err)
					}
				}()
			}

		}
		wg.Wait()
		account1, err := user.GetBalance(1)
		if err != nil {
			t.Error(err)
			return
		}
		account2, err := user.GetBalance(2)
		if err != nil {
			t.Error(err)
			return
		}
		account1BalanceWant := 250.0
		if account1.Balance != account1BalanceWant {
			t.Errorf("Transfer() failed. Debtor balance want: %v, got %v", account1BalanceWant, account1.Balance)
		}
		account2BalanceWant := 750.0
		if account2.Balance != account2BalanceWant {
			t.Errorf("Transfer() failed. Beneciary balance want: %v, got %v", account2BalanceWant, account2.Balance)
		}
	})
	t.Run("Should not be able to make a transfer because something went wrong while doing it and the balance of the debtor should be restored", func(t *testing.T) {
		DB.Exec("DELETE FROM users;")
		defer DB.Exec("DELETE FROM users;")
		_, err = DB.Exec("INSERT INTO users (id, username, balance) VALUES (1, 'first_user', 500);")
		if err != nil {
			t.Error(err)
			return
		}
		_, err = DB.Exec("INSERT INTO users (id, username, balance) VALUES (2, 'second_user', 100);")
		if err != nil {
			t.Error(err)
			return
		}
		uowFakeRepo := uow.NewUowImpl(DB, repository.NewUserRepositoryFakeTest())
		user := service.NewUser(DB, uowFakeRepo)
		err := user.Transfer(&service.TransferInput{
			Amount:        100,
			DebtorID:      1,
			BeneficiaryID: 2,
		})
		outputDebtor, err := user.GetBalance(1)
		if err != nil {
			t.Error(err)
			return
		}
		outputBeneficiary, err := user.GetBalance(2)
		if err != nil {
			t.Error(err)
			return
		}
		outputDebtorWant := 500.0
		if outputDebtor.Balance != outputDebtorWant {
			t.Errorf("Transfer() failed. want: %v, got %v", outputDebtorWant, outputDebtor.Balance)
		}
		outputBeneficiaryWant := 100.0
		if outputBeneficiary.Balance != outputBeneficiaryWant {
			t.Errorf("Transfer() failed. want: %v, got %v", outputBeneficiaryWant, outputBeneficiary.Balance)
		}
	})
}

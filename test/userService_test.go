package test

import (
	"errors"
	"money-transfer-api/infra/database"
	"money-transfer-api/repository"
	"money-transfer-api/service"
	"money-transfer-api/uow"
	"sync"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func TestUserService(t *testing.T) {
	DB, err := database.Open()
	if err != nil {
		panic(err)
	}
	defer DB.Close()
	uow := uow.NewUowImpl(DB)
	t.Run("Should get the user balance by its ID", func(t *testing.T) {
		DB.Exec("DELETE FROM users;")
		defer DB.Exec("DELETE FROM users;")
		_, err = DB.Exec("INSERT INTO users (id, username, balance) VALUES (1, 'first_user', 2000);")
		if err != nil {
			t.Error(err)
			return
		}
		user := service.NewUser(DB, uow, repository.USER_REPOSITORY_POSTGRES)
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
		user := service.NewUser(DB, uow, repository.USER_REPOSITORY_POSTGRES)
		var wg sync.WaitGroup
		totalRequests := 5
		wg.Add(totalRequests)
		for i := 0; i < totalRequests; i++ {
			time.Sleep(time.Millisecond * 10)
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
		user := service.NewUser(DB, uow, repository.USER_REPOSITORY_POSTGRES)
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
		user := service.NewUser(DB, uow, repository.FAKE_USER_REPOSITORY_DEPOSIT)
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
		user := service.NewUser(DB, uow, repository.USER_REPOSITORY_POSTGRES)
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
}

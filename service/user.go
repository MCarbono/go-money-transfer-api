package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"money-transfer-api/uow"
)

type User struct {
	DB  *sql.DB
	uow uow.Uow
}

func NewUser(db *sql.DB, uow uow.Uow) *User {
	return &User{
		DB:  db,
		uow: uow,
	}
}

func (u *User) Transfer(input *TransferInput) (err error) {
	err = u.uow.Do(context.Background(), func(uow *uow.UowImpl) error {
		var getBalanceDTOOutput GetBalanceDTOOutput
		err = uow.Tx.QueryRow(`SELECT balance FROM users WHERE id = $1`, input.DebtorID).Scan(&getBalanceDTOOutput.Balance)
		if err != nil {
			return fmt.Errorf("error trying to get balance from user with ID %v. %w", input.DebtorID, err)
		}
		if getBalanceDTOOutput.Balance < input.Amount {
			return errors.New("insufficient funds")
		}
		newBalance := getBalanceDTOOutput.Balance - input.Amount
		_, err = uow.Tx.Exec("UPDATE users SET balance = $2 WHERE id = $1", input.DebtorID, newBalance)
		if err != nil {
			return err
		}
		var getBalanceDTOOutputBene GetBalanceDTOOutput
		err = uow.Tx.QueryRow(`SELECT balance FROM users WHERE id = $1`, input.BeneficiaryID).Scan(&getBalanceDTOOutputBene.Balance)
		if err != nil {
			return fmt.Errorf("error trying to get balance from user with ID %v. %w", input.BeneficiaryID, err)
		}
		b := getBalanceDTOOutputBene.Balance + input.Amount
		_, err = uow.Tx.Exec("UPDATE users SET balance = $2 WHERE id = $1", input.BeneficiaryID, b)
		return err
	})
	return
}

type TransferInput struct {
	Amount        float64
	DebtorID      int
	BeneficiaryID int
}

func (u *User) GetBalance(userID int) (GetBalanceDTOOutput, error) {
	var getBalanceDTOOutput GetBalanceDTOOutput
	err := u.DB.QueryRow(`SELECT balance FROM users WHERE id = $1`, userID).Scan(&getBalanceDTOOutput.Balance)
	if err != nil {
		return GetBalanceDTOOutput{}, fmt.Errorf("error trying to get balance from user with ID %v. %w", userID, err)
	}
	return getBalanceDTOOutput, nil
}

type GetBalanceDTOOutput struct {
	Balance float64 `json:"balance"`
}

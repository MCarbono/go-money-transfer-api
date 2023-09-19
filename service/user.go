package service

import (
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

func (u *User) Transfer(input *TransferInput) error {
	outputBalanceDebtor, err := u.GetBalance(input.BeneficiaryID)
	if err != nil {
		return fmt.Errorf("error trying to get balance from user with ID %v. %w", input.DebtorID, err)
	}
	if outputBalanceDebtor.Balance < input.Amount {
		return errors.New("insufficient funds")
	}
	newBalance := outputBalanceDebtor.Balance - input.Amount
	_, err = u.DB.Exec("UPDATE users SET balance = $2 WHERE id = $1", input.DebtorID, newBalance)
	if err != nil {
		return fmt.Errorf("error trying to withdraw: %w", err)
	}

	return nil
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

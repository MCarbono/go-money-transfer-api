package service

import (
	"context"
	"database/sql"
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
	err = u.uow.Do(context.Background(), func(tx *sql.Tx) error {
		ctx := context.Background()
		repo := u.uow.GetUserRepository(ctx, tx)
		debtorUser, err := repo.FindUserTx(ctx, input.DebtorID)
		if err != nil {
			return err
		}
		err = debtorUser.Withdraw(input.Amount)
		if err != nil {
			return err
		}
		beneficiaryUser, err := repo.FindUserTx(ctx, input.BeneficiaryID)
		if err != nil {
			return err
		}
		beneficiaryUser.Deposit(input.Amount)
		if beneficiaryUser.ID > debtorUser.ID {
			err = repo.Update(ctx, debtorUser)
			if err != nil {
				return err
			}
			err = repo.Update(ctx, beneficiaryUser)
			return err
		} else {
			err = repo.Update(ctx, beneficiaryUser)
			if err != nil {
				return err
			}
			err = repo.Update(ctx, debtorUser)
			return err
		}
	})
	return
}

type TransferInput struct {
	Amount        float64 `json:"amount"`
	DebtorID      int     `json:"debtor_id"`
	BeneficiaryID int     `json:"beneficiary_id"`
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

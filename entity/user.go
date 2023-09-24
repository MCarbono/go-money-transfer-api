package entity

import "errors"

type User struct {
	ID       int
	Username string
	Balance  float64
}

func (u *User) Withdraw(amount float64) error {
	if u.Balance < amount {
		return errors.New("insufficient funds")
	}
	u.Balance -= amount
	return nil
}

func (u *User) Deposit(amount float64) {
	u.Balance += amount
}

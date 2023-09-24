package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"money-transfer-api/entity"
)

const (
	USER_REPOSITORY_POSTGRES     = "postgres"
	FAKE_USER_REPOSITORY_DEPOSIT = "fake"
)

type UserRepository interface {
	FindUserTx(ctx context.Context, ID int) (*entity.User, error)
	Withdraw(ctx context.Context, user *entity.User) error
	Deposit(ctx context.Context, user *entity.User) error
}

type UserRepositoryPostgres struct {
	DB *sql.DB
	Tx *sql.Tx
}

func NewUserRepositoryPostgres(DB *sql.DB, TX *sql.Tx) *UserRepositoryPostgres {
	return &UserRepositoryPostgres{
		DB: DB,
		Tx: TX,
	}
}

func (r *UserRepositoryPostgres) FindUserTx(ctx context.Context, ID int) (*entity.User, error) {
	var user entity.User
	err := r.Tx.QueryRow(`SELECT * FROM users WHERE id = $1`, ID).Scan(&user.ID, &user.Username, &user.Balance)
	if err != nil {
		return nil, fmt.Errorf("error trying to get user with ID %v. %w", ID, err)
	}
	return &user, nil
}

func (r *UserRepositoryPostgres) Withdraw(ctx context.Context, user *entity.User) error {
	_, err := r.Tx.Exec("UPDATE users SET balance = $2 WHERE id = $1", user.ID, user.Balance)
	if err != nil {
		return fmt.Errorf("error trying to withdraw. %w", err)
	}
	return nil
}

func (r *UserRepositoryPostgres) Deposit(ctx context.Context, user *entity.User) error {
	_, err := r.Tx.Exec("UPDATE users SET balance = $2 WHERE id = $1", user.ID, user.Balance)
	if err != nil {
		return fmt.Errorf("error trying to deposit balance. %w", err)
	}
	return err
}

type UserRepositoryFakeTest struct {
	DB *sql.DB
	Tx *sql.Tx
}

func NewUserRepositoryFakeTest(DB *sql.DB, TX *sql.Tx) *UserRepositoryFakeTest {
	return &UserRepositoryFakeTest{
		DB: DB,
		Tx: TX,
	}
}

func (r *UserRepositoryFakeTest) FindUserTx(ctx context.Context, ID int) (*entity.User, error) {
	var user entity.User
	err := r.Tx.QueryRow(`SELECT * FROM users WHERE id = $1`, ID).Scan(&user.ID, &user.Username, &user.Balance)
	if err != nil {
		return nil, fmt.Errorf("error trying to get user with ID %v. %w", ID, err)
	}
	return &user, nil
}

func (r *UserRepositoryFakeTest) Withdraw(ctx context.Context, user *entity.User) error {
	_, err := r.Tx.Exec("UPDATE users SET balance = $2 WHERE id = $1", user.ID, user.Balance)
	if err != nil {
		return fmt.Errorf("error trying to withdraw. %w", err)
	}
	return nil
}

func (r *UserRepositoryFakeTest) Deposit(ctx context.Context, user *entity.User) error {
	return errors.New("Test error")
}

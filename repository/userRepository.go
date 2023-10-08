package repository

import (
	"context"
	"database/sql"
	"fmt"
	"money-transfer-api/entity"
)

type UserRepository interface {
	FindUserTx(ctx context.Context, ID int) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
}

type UserRepositoryPostgres struct {
	Tx *sql.Tx
}

func NewUserRepositoryPostgres(TX *sql.Tx) *UserRepositoryPostgres {
	return &UserRepositoryPostgres{
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

func (r *UserRepositoryPostgres) Update(ctx context.Context, user *entity.User) error {
	_, err := r.Tx.Exec("UPDATE users SET balance = $2, username = $3 WHERE id = $1", user.ID, user.Balance, user.Username)
	if err != nil {
		return fmt.Errorf("error trying to update. %w", err)
	}
	return nil
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

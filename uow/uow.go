package uow

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"money-transfer-api/repository"
)

type Uow interface {
	Do(ctx context.Context, fn func(uow *UowImpl) error) error
	CommitOrRollback() error
	Rollback() error
	GetUserRepository(ctx context.Context, implementationName string) (repository.UserRepository, error)
}

type UowImpl struct {
	Db *sql.DB
	Tx *sql.Tx
}

func NewUowImpl(db *sql.DB) *UowImpl {
	return &UowImpl{
		Db: db,
	}
}

func (u *UowImpl) GetUserRepository(ctx context.Context, implementationName string) (repository.UserRepository, error) {
	var repo repository.UserRepository
	if implementationName == repository.USER_REPOSITORY_POSTGRES {
		repo = repository.NewUserRepositoryPostgres(u.Db, u.Tx)
	}
	if implementationName == repository.FAKE_USER_REPOSITORY_DEPOSIT {
		repo = repository.NewUserRepositoryFakeTest(u.Db, u.Tx)
	}
	return repo, nil
}

func (u *UowImpl) Do(ctx context.Context, fn func(uow *UowImpl) error) error {
	if u.Tx != nil {
		return errors.New("transaction already started")
	}
	tx, err := u.Db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	u.Tx = tx
	err = fn(u)
	if err != nil {
		errRb := u.Rollback()
		if errRb != nil {
			return fmt.Errorf("orignal error: %s, rollback error: %s", err.Error(), errRb.Error())
		}
		return err
	}
	return u.CommitOrRollback()
}

func (u *UowImpl) CommitOrRollback() error {
	defer func() {
		u.Tx = nil
	}()
	err := u.Tx.Commit()
	if err != nil {
		errRb := u.Rollback()
		if errRb != nil {
			return fmt.Errorf("orignal error: %s, rollback error: %s", err.Error(), errRb.Error())
		}
		return err
	}
	u.Tx = nil
	return nil
}

func (u *UowImpl) Rollback() error {
	defer func() {
		u.Tx = nil
	}()
	if u.Tx == nil {
		return errors.New("no transaction to rollback")
	}
	err := u.Tx.Rollback()
	if err != nil {
		return err
	}
	u.Tx = nil
	return nil
}

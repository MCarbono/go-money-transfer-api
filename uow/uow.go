package uow

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"money-transfer-api/repository"
)

type RepositoryFactory func() interface{}

type Uow interface {
	Register(name string, fc RepositoryFactory)
	Do(ctx context.Context, fn func(uow *UowImpl) error) error
	CommitOrRollback() error
	Rollback() error
	GetUserRepository(ctx context.Context, implementationName string) (repository.UserRepository, error)
}

type UowImpl struct {
	Db           *sql.DB
	Tx           *sql.Tx
	Repositories map[string]RepositoryFactory
}

func (u *UowImpl) Register(name string, fc RepositoryFactory) {
	u.Repositories[name] = fc
}

func NewUowImpl(db *sql.DB) *UowImpl {
	return &UowImpl{
		Db:           db,
		Repositories: make(map[string]RepositoryFactory),
	}
}

func (u *UowImpl) GetUserRepository(ctx context.Context, implementationName string) (repository.UserRepository, error) {
	repo := u.Repositories[implementationName]()
	return repo.(repository.UserRepository), nil
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

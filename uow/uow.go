package uow

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"money-transfer-api/repository"
)

var isolationSerializableErr = errors.New("ERROR: could not serialize access due to concurrent update (SQLSTATE 40001)")

type Uow interface {
	Do(ctx context.Context, fn func(tx *sql.Tx) error) error
	GetUserRepository(ctx context.Context, tx *sql.Tx) repository.UserRepository
}

type UowImpl struct {
	Db             *sql.DB
	totalRetries   int
	userRepository repository.UserRepository
}

func NewUowImpl(db *sql.DB, userRepository repository.UserRepository) *UowImpl {
	return &UowImpl{
		Db:             db,
		totalRetries:   7,
		userRepository: userRepository,
	}
}

func (u *UowImpl) GetUserRepository(ctx context.Context, tx *sql.Tx) repository.UserRepository {
	return u.userRepository.Clone(tx)
}

func (u *UowImpl) Do(ctx context.Context, fn func(tx *sql.Tx) error) error {
	var err error
	var tx *sql.Tx
	var retries = 0
	for retries < u.totalRetries {
		tx, err = u.Db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
		if err != nil {
			return err
		}
		err = fn(tx)
		if err == nil {
			break
		}
		if err != nil {
			if errors.Unwrap(err) != nil {
				if errors.Unwrap(err).Error() == isolationSerializableErr.Error() {
					errRb := tx.Rollback()
					if errRb != nil {
						return fmt.Errorf("original error: %s, rollback error: %s", err.Error(), errRb.Error())
					}
					retries++
					continue
				}
			}
			break
		}
	}
	if err != nil {
		errRb := tx.Rollback()
		if errRb != nil {
			return fmt.Errorf("original error: %s, rollback error: %s", err.Error(), errRb.Error())
		}
		return err
	}
	err = tx.Commit()
	if err != nil {
		errRb := tx.Rollback()
		if errRb != nil {
			return fmt.Errorf("commit error: %s, rollback error: %s", err.Error(), errRb.Error())
		}
	}
	return nil
}

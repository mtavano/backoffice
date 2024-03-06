package service

import (
	"context"
	"database/sql"

	"log"

	"github.com/darchlabs/backoffice/internal/storage"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func (ng *Engine) WithTransaction(db storage.Transaction, fn func(txx storage.Transaction) error) (err error) {
	ctx := context.Background()
	var tx *sqlx.Tx

	defer func() {
		if err != nil {
			if txErr := tx.Rollback(); txErr != nil {
				err = errors.Wrapf(txErr, "reward: Worker.WithTransaction error %s", err)
				log.Printf("[pkg: service] Enfine.WithTransaction tx.Rollback [error: '%s']", err)
			}
		}
	}()

	switch store := db.(type) {
	case *sqlx.DB:
		tx, err = store.BeginTxx(ctx, &sql.TxOptions{
			Isolation: sql.LevelSerializable,
		})
		if err != nil {
			err = errors.Wrap(err, "reward: Worker.WithTransaction error")
			log.Printf("[pkg: service] Enfine.WithTransaction [error: '%s']", err)
			return err
		}
	case *sqlx.Tx:
		tx = store

	default:
		Tx, err := ng.database.BeginTx(ctx)
		if err != nil {
			err = errors.Wrap(err, "reward: Worker.WithTransaction w.db.BeginTx error")
			log.Printf("[pkg: service] Enfine.WithTransaction [error: '%s']", err)
			return err
		}

		tx = Tx.(*sqlx.Tx)
	}

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit()
}

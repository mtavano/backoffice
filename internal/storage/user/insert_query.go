package user

import (
	"github.com/darchlabs/backoffice/internal/storage"
	"github.com/pkg/errors"
)

func InsertQuery(tx storage.QueryContext, record *Record) error {
	_, err := tx.Exec(`
		INSERT INTO users (
			id,
			email,
			hashed_password,
			verified,
			created_at
		) VALUES ($1, $2, $3, $4, $5);`,
		record.ID,
		record.Email,
		record.HashedPassword,
		record.Verified,
		record.CreatedAt,
	)
	if err != nil {
		return errors.Wrap(err, "user: InsertQuery tx.Exec error")
	}
	return nil
}

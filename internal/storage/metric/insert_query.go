package metric

import (
	"github.com/darchlabs/backoffice/internal/storage"
	"github.com/pkg/errors"
)

func InsertQuery(tx storage.QueryContext, record *Record) error {
	_, err := tx.Exec(`
		INSERT INTO metrics (
			user_id,
			profile_id,
			prints,
			created_at
		) VALUES ($1, $2, $3, $4);`,
		record.UserID,
		record.ProfileID,
		record.Prints,
		record.CreatedAt,
	)
	if err != nil {
		return errors.Wrap(err, "user: InsertQuery tx.Exec error")
	}
	return nil
}

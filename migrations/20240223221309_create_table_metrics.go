package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upCreateTableMetrics, downCreateTableMetrics)
}

func upCreateTableMetrics(tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	_, err := tx.Exec(`
		CREATE TABLE metrics (
			user_id TEXT REFERENCES users(id),
			profile_id TEXT REFERENCES profiles(short_id),
			prints integer NOT NULL default 0,

			created_at TIMESTAMPTZ NOT NULL,
			updated_at TIMESTAMPTZ,
			deleted_at TIMESTAMPTZ
		);
	`)
	return err
}

func downCreateTableMetrics(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}

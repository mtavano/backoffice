package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upCreateTableUsers, downCreateTableUsers)
}

func upCreateTableUsers(tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	_, err := tx.Exec(`
		CREATE TABLE users (
			id TEXT PRIMARY KEY NOT NULL,
			email TEXT NOT NULL UNIQUE,
			nickname TEXT UNIQUE NOT NULL,
			hashed_password TEXT NOT NULL,
			verified BOOLEAN DEFAULT false,
			created_at TIMESTAMPTZ NOT NULL,
			updated_at TIMESTAMPTZ,
			deleted_at TIMESTAMPTZ
		);`)
	if err != nil {
		return err
	}

	return nil
}

func downCreateTableUsers(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	_, err := tx.Exec(`DROP TABLE users;`)
	if err != nil {
		return err
	}

	return nil
}

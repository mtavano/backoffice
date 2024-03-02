package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upCreateTableCards, downCreateTableCards)
}

func upCreateTableCards(tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	_, err := tx.Exec(`
		CREATE TABLE cards (
			short_id TEXT NOT NULL PRIMARY KEY,
			created_at TIMESTAMPTZ default now(),
			status TEXT DEFAULT 'free'
		);`)
	return err
}

func downCreateTableCards(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	_, err := tx.Exec(`DROP TABLE cards;`)
	return err
}

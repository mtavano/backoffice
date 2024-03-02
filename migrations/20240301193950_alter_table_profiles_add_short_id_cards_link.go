package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upAlterTableProfilesAddShortIdCardsLink, downAlterTableProfilesAddShortIdCardsLink)
}

func upAlterTableProfilesAddShortIdCardsLink(tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	_, err := tx.Exec(`
		ALTER TABLE profiles
		ADD CONSTRAINT fk_profiles_cards
		FOREIGN KEY (short_id) REFERENCES cards(short_id)
		ON DELETE CASCADE;
	`)
	return err
}

func downAlterTableProfilesAddShortIdCardsLink(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	_, err := tx.Exec(`
		ALTER TABLE profiles
		DROP COLUMN IF EXISTS card_short_id;`)
	return err
}

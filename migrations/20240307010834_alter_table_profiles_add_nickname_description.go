package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upAlterTableProfilesAddNicknameDescription, downAlterTableProfilesAddNicknameDescription)
}

func upAlterTableProfilesAddNicknameDescription(tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	// drop column from users
	_, err := tx.Exec(`
		ALTER TABLE users
		DROP COLUMN nickname;
	`)
	if err != nil {
		return err
	}
	// add columns to profiles
	_, err = tx.Exec(`
		ALTER TABLE profiles
		ADD COLUMN description TEXT DEFAULT '',
		ADD COLUMN nickname TEXT UNIQUE;
	`)

	return nil
}

func downAlterTableProfilesAddNicknameDescription(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}

package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upCreateTableProfile, downCreateTableProfile)
}

func upCreateTableProfile(tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	_, err := tx.Exec(`
		CREATE TABLE profiles (
			short_id TEXT UNIQUE NOT NULL,
			user_id TEXT PRIMARY KEY REFERENCES users(id),
			linkedin TEXT DEFAULT '',
			email TEXT DEFAULT '',
			whatsapp TEXT DEFAULT '',
			medium TEXT DEFAULT '',
			website TEXT DEFAULT '',
			twitter_x TEXT DEFAULT '',
			image TEXT DEFAULT '',

			created_at TIMESTAMPTZ NOT NULL,
			updated_at TIMESTAMPTZ,
			deleted_at TIMESTAMPTZ
		);`)
	if err != nil {
		return err
	}
	return nil
}

func downCreateTableProfile(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	_, err := tx.Exec("DROP TABLE profile;")
	return err
}

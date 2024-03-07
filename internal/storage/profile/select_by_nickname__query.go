package profile

import (
	"database/sql"

	"github.com/darchlabs/backoffice/internal/storage"
	"github.com/pkg/errors"
)

func SelectByNicknameQuery(tx storage.Transaction, nn string) (*Record, error) {
	var record Record
	err := tx.Get(&record, `
		SELECT *
		FROM profiles
		WHERE nickname = $1;`,
		nn,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNoProfile
	}
	if err != nil {
		return nil, errors.Wrap(err, "profile: SelectByNickname tx.Get error")
	}
	return &record, nil
}

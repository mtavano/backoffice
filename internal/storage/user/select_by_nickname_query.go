package user

import (
	"database/sql"

	"github.com/darchlabs/backoffice/internal/storage"
	"github.com/pkg/errors"
)

func SelectByNicknameQuery(tx storage.Transaction, nn string) (*Record, error) {
	var record Record
	err := tx.Get(&record, `
		SELECT *
		FROM users
		WHERE nickname = $1;`,
		nn,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, errors.Wrap(err, "user: SelectByNickname tx.Get error")
	}
	return &record, nil
}

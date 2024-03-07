package profile

import (
	"database/sql"

	"github.com/darchlabs/backoffice/internal/storage"
	"github.com/pkg/errors"
)

type SelectFilters struct {
	ShortID  string
	Nickname string
}

func SelectQuery(tx storage.Transaction, filters *SelectFilters) (*Record, error) {
	var r *Record
	if filters.ShortID != "" {
		rr, err := selectByShortID(tx, filters.ShortID)
		if err != nil {
			return nil, err
		}
		return rr, nil
	}

	r, err := selectByNickname(tx, filters.Nickname)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func selectByShortID(tx storage.Transaction, shortID string) (*Record, error) {
	var record Record
	err := tx.Get(&record, `
		SELECT *
		FROM profiles
		WHERE short_id = $1;
	`, shortID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNoProfile
	}
	if err != nil {
		return nil, errors.Wrap(err, "profile: SelectByNickname tx.Get error")
	}

	return &record, nil
}

func selectByNickname(tx storage.Transaction, nickname string) (*Record, error) {
	var record Record
	err := tx.Get(&record, `
	SELECT *
	FROM profiles p
	WHERE
    p.nickname = $1;`,
		nickname,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNoProfile
	}
	if err != nil {
		return nil, errors.Wrap(err, "profile: SelectByNickname tx.Get error")
	}

	return &record, nil
}

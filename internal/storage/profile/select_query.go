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
	if filters.ShortID != "" {
		return selectByShortID(tx, filters.ShortID)
	}
	if filters.Nickname != "" {
		return selectByNickname(tx, filters.Nickname)
	}

	return nil, sql.ErrNoRows
}

func selectByShortID(tx storage.Transaction, shortID string) (*Record, error) {
	var record Record
	err := tx.Get(&record, `SELECT * FROM profiles WHERE short_id = $1;`, shortID)
	if err != nil {
		return nil, errors.Wrap(err, "profile: SelectByNickname tx.Get error")
	}

	return &record, nil
}

func selectByNickname(tx storage.Transaction, nickname string) (*Record, error) {
	var record Record
	err := tx.Get(&record, `
	SELECT
    u.email,
    p.short_id,
    p.linkedin,
    p.email,
    p.whatsapp,
    p.medium,
    p.website,
    p.twitter_x,
    p.image
	FROM
    users u
	JOIN
    profiles p ON u.id = p.user_id
	WHERE
    u.nickname = $1;`,
		nickname,
	)
	if err != nil {
		return nil, errors.Wrap(err, "profile: SelectByNickname tx.Get error")
	}
	return &record, nil
}

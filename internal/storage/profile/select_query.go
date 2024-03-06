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

	return selectByNickname(tx, filters.Nickname)
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

	err = tx.Get(&record.Nickname, `
		SELECT nickname
		FROM users
		WHERE id = $1;`,
		record.UserID,
	)
	if err != nil {
		return nil, errors.Wrap(err, "profile: SelectByNickname tx.Get error")
	}

	return &record, nil
}

func selectByNickname(tx storage.Transaction, nickname string) (*Record, error) {
	var record Record
	err := tx.Get(&record, `
	SELECT
    p.email,
    p.short_id,
    p.linkedin,
    p.email,
    p.whatsapp,
    p.medium,
    p.website,
    p.twitter_x,
    p.image,
		p.created_at,
		p.updated_at
	FROM
    users u
	JOIN
    profiles p ON u.id = p.user_id
	WHERE
    u.nickname = $1;`,
		nickname,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNoProfile
	}
	if err != nil {
		return nil, errors.Wrap(err, "profile: SelectByNickname tx.Get error")
	}

	record.Nickname = nickname
	return &record, nil
}

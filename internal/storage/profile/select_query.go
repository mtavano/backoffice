package profile

import (
	"database/sql"
	"log"

	"github.com/darchlabs/backoffice/internal/storage"
	"github.com/pkg/errors"
)

type SelectFilters struct {
	ShortID  string
	Nickname string
	UserID   string
}

func SelectQuery(tx storage.Transaction, f *SelectFilters) (r *Record, err error) {
	switch true {
	case f.ShortID != "":
		r, err = selectByShortID(tx, f.ShortID)
		break

	case f.Nickname != "":
		r, err = selectByNickname(tx, f.Nickname)
		break

	case f.UserID != "":
		r, err = selectByUserID(tx, f.UserID)
		break

	default:
		return nil, ErrNoProfile
	}

	return r, err
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
		return nil, errors.Wrap(err, "profile: selectByNickname tx.Get error")
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
		return nil, errors.Wrap(err, "profile: selectByNickname tx.Get error")
	}

	return &record, nil
}

func selectByUserID(tx storage.Transaction, userID string) (*Record, error) {
	var record Record
	err := tx.Get(&record, `
	SELECT *
	FROM profiles p
	WHERE
    p.user_id = $1;`,
		userID,
	)
	if errors.Is(err, sql.ErrNoRows) {
		log.Println("errrorrrrrr")
		return nil, ErrNoProfile
	}
	if err != nil {
		return nil, errors.Wrap(err, "profile: selectByUserID tx.Get error")
	}

	return &record, nil
}

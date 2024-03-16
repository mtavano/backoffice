package profile

import (
	"time"

	"github.com/darchlabs/backoffice/internal/storage"
	"github.com/pkg/errors"
)

type UpdateNicknameInput struct {
	UserID   string
	ShortID  string
	Nickname *string

	Time time.Time
}

func UpdateNicknameQuery(tx storage.Transaction, input *UpdateNicknameInput) (*Record, error) {
	// mark card as claimed
	var record Record
	err := tx.Get(&record, `
		UPDATE profiles
		SET nickname = $1,
		updated_at = $4
		WHERE profiles.short_id = $2 
		AND profiles.user_id = $3 -- Match both user_id and short_id
		RETURNING *;
		;`,
		input.Nickname,
		input.ShortID,
		input.UserID,
		input.Time,
	)
	if err != nil {
		return nil, errors.Wrap(err, "profile: UpsertQuery tx.Exec error")
	}

	return &record, nil
}

package profile

import (
	"time"

	"github.com/darchlabs/backoffice/internal/storage"
	"github.com/pkg/errors"
)

type UpsertProfileInput struct {
	ShortID string
	UserID  string

	// Social network links
	Linkedin    *string
	Email       *string
	Whatsapp    *string
	Medium      *string
	TwitterX    *string
	Website     *string
	Description *string

	Time time.Time
}

func UpsertQuery(tx storage.Transaction, input *UpsertProfileInput) (*Record, error) {
	// mark card as claimed
	var record Record
	err := tx.Get(&record, `
		INSERT INTO profiles (
				user_id,    -- 1
				short_id,   -- 2 
				linkedin,   -- 3
				email,      -- 4
				whatsapp,   -- 5
				medium,     -- 6
				website,    -- 7
				created_at, -- 8
				twitter_x,   -- 9
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (short_id) -- Add short_id to the conflict target
		DO UPDATE
		SET 
				linkedin =$3,
				email =$4,
				whatsapp =$5,
				medium =$6,
				website =$7,
				twitter_x =$9,
				updated_at = $8
		WHERE profiles.short_id = $2 
		AND profiles.user_id = $1 -- Match both user_id and short_id
		RETURNING *;
		;`,
		input.UserID,
		input.ShortID,
		input.Linkedin,
		input.Email,
		input.Whatsapp,
		input.Medium,
		input.Website,
		input.Time,
		input.TwitterX,
	)
	if err != nil {
		return nil, errors.Wrap(err, "profile: UpsertQuery tx.Exec error")
	}

	return &record, nil
}

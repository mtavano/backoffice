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
	Nickname    *string

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
				description, -- 10
				nickname -- 11
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (short_id) -- Add short_id to the conflict target
		DO UPDATE
		SET 
				linkedin = COALESCE($3, profiles.linkedin),
				email = COALESCE($4, profiles.email),
				whatsapp = COALESCE($5, profiles.whatsapp),
				medium = COALESCE($6, profiles.medium),
				website = COALESCE($7, profiles.website),
				twitter_x = COALESCE($9, profiles.twitter_x),
				description = COALESCE($10, profiles.description),
				nickname = COALESCE($11, profiles.nickname),
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
		input.Description,
		input.Nickname,
	)
	if err != nil {
		return nil, errors.Wrap(err, "profile: UpsertQuery tx.Exec error")
	}

	return &record, nil
}

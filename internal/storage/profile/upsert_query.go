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
	Linkedin *string
	Email    *string
	Whatsapp *string
	Medium   *string
	TwitterX *string
	Website  *string

	Time time.Time
}

func UpsertQuery(tx storage.Transaction, input *UpsertProfileInput) (*Record, error) {
	// mark card as claimed
	var record Record
	err := tx.Get(&record, `
		INSERT INTO profiles (
			user_id, -- 1
			short_id, -- 2 
			linkedin,-- 3
			email, -- 4
			whatsapp, -- 5
			medium, -- 6
			website, -- 7
			created_at, -- 8
			twitter_x -- 9
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (user_id)
		DO UPDATE
		SET 
			linkedin = COALESCE($3, profiles.linkedin),
			email = COALESCE($4, profiles.email),
			whatsapp = COALESCE($5, profiles.whatsapp),
			medium = COALESCE($6, profiles.medium),
			website = COALESCE($7, profiles.website),
			twitter_x = COALESCE($9, profiles.twitter_x),
			updated_at = $8
		WHERE profiles.short_id = $2
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

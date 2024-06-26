package profile

import (
	"errors"
	"time"
)

var ErrNoProfile = errors.New("profile: no profile erro")
var ErrInvalidFilters = errors.New("profile: invalid filters for query")

type Record struct {
	UserID  string `db:"user_id"`
	ShortID string `db:"short_id"`

	// Social network links
	Linkedin    *string `db:"linkedin"`
	Email       *string `db:"email"`
	Whatsapp    *string `db:"whatsapp"`
	Medium      *string `db:"medium"`
	TwitterX    *string `db:"twitter_x"`
	Website     *string `db:"website"`
	Description *string `db:"description"`
	Nickname    *string `db:"nickname"`

	// Non available fort the moment
	Image *string `db:"image"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

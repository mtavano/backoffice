package metric

import "time"

type Record struct {
	UserID    string `db:"user_id"`
	ProfileID string `db:"profile_id"`
	Prints    uint64 `db:"prints"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

package cards

import "time"

const (
	StatusFree    = "free"
	StatusClaimed = "claimed"
)

type Record struct {
	ShortID   string    `db:"short_id"`
	Status    string    `db:"status"`
	CreatedAt time.Time `db:"created_at"`
}

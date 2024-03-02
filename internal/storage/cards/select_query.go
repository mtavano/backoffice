package cards

import (
	"database/sql"

	"github.com/darchlabs/backoffice/internal/storage"
	"github.com/pkg/errors"
)

var ErrNoCard = errors.New("no card error")

func SelectQuery(tx storage.Transaction, shortID string) (*Record, error) {
	var record Record
	err := tx.Get(&record, `
		SELECT *
		FROM cards
		WHERE short_id = $1;`,
		shortID,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNoCard
	}
	if err != nil {
		return nil, errors.Wrap(err, "cards: SelectQuery tx.Get error")
	}

	return &record, nil
}

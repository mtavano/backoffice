package cards

import (
	"github.com/darchlabs/backoffice/internal/storage"
	"github.com/pkg/errors"
)

func UpdateQuery(qc storage.QueryContext, shortID string) error {
	_, err := qc.Exec(`
		UPDATE cards
		SET status = 'claimed'
		WHERE short_id = $1;`,
		shortID,
	)
	if err != nil {
		return errors.Wrap(err, "cards: UpdateQuery qc.Exec error")
	}
	return nil
}

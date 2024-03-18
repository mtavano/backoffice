package cards

import (
	"fmt"
	"strings"

	"github.com/darchlabs/backoffice/internal/storage"
	"github.com/pkg/errors"
)

func InsertQuery(tx storage.Transaction, ids []string) error {
	values := make([]string, 0)
	for _, id := range ids {
		values = append(values, fmt.Sprintf("('%s')", id))
	}

	query := fmt.Sprintf(`
		INSERT INTO cards (short_id)
		VALUES %s;`,
		strings.Join(values, ","),
	)

	_, err := tx.Exec(query)
	if err != nil {
		return errors.Wrap(err, "admin: InsertQuery tx.Exec")
	}

	return nil
}

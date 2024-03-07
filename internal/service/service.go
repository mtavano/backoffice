package service

import (
	"github.com/darchlabs/backoffice/internal/storage"
	"github.com/darchlabs/backoffice/internal/storage/profile"
	"github.com/pkg/errors"
)

type Engine struct {
	database *storage.SQLStore
}

func (ng *Engine) AtomicUpsertCard(input *profile.UpsertProfileInput) (*profile.Record, error) {
	var r *profile.Record
	err := ng.WithTransaction(ng.database, func(txx storage.Transaction) error {
		// check wether card exist or not
		// check if card is claimed by give user id
		// upsert
		return nil
	})

	if err != nil {
		return nil, errors.Wrap(err, "service: Engine.AtomicUpsertCard ng.WithTransaction error")
	}
	return r, nil

}

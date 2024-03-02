package admin

import (
	"fmt"

	"github.com/darchlabs/backoffice/internal/api/context"
	v1 "github.com/darchlabs/backoffice/internal/api/v1"
	cardsdb "github.com/darchlabs/backoffice/internal/storage/cards"
	"github.com/darchlabs/backoffice/pkg/middleware"
)

func Route(basePath string, ctx *context.Ctx) {
	adm := middleware.NewAdminKey(ctx.App.Config.AdminKey)

	// handlers
	postCardsHandler := &PostCardsHandler{
		insertCardsQuery: cardsdb.InsertQuery,
	}

	ctx.Server.Post(
		fmt.Sprintf("%s/cards", basePath),
		adm.Middleware,
		v1.HandleFunc(ctx, postCardsHandler.Invoke),
	)

	return
}

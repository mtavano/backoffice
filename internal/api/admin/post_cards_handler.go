package admin

import (
	"github.com/darchlabs/backoffice/internal/api/context"
	"github.com/darchlabs/backoffice/internal/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

type insertCardsQuery func(storage.Transaction, []string) error

type PostCardsHandler struct {
	insertCardsQuery insertCardsQuery
}

type postCardsHandlerRequest struct {
	TotalCards int `json:"totalCards"`
}

type postCardsHandlerResponse struct {
	ShortIDs []string `json:"shortIds"`
}

func (h *PostCardsHandler) Invoke(ctx *context.Ctx, c *fiber.Ctx) (interface{}, int, error) {
	var req postCardsHandlerRequest
	err := c.BodyParser(&req)
	if err != nil {
		return nil, fiber.StatusInternalServerError, errors.Wrap(
			err,
			"admin: PostCardsHandler.Invoke c.BodyParser error",
		)
	}

	ids := make([]string, 0)
	for i := 0; i < req.TotalCards; i++ {
		sid, err := ctx.ShortIDGenerator()
		if err != nil {
			return nil, fiber.StatusInternalServerError, errors.Wrap(
				err, "admin: PostCardsHandler.Invoke ctx.ShortIDGenerator error",
			)
		}

		ids = append(ids, sid)
	}

	err = h.insertCardsQuery(ctx.App.SqlStore, ids)
	if err != nil {
		return nil, fiber.StatusInternalServerError, errors.Wrap(
			err, "admin: PostCardsHandler.Invoke h.insertCardsQuery error",
		)
	}

	return &postCardsHandlerResponse{ShortIDs: ids}, fiber.StatusCreated, nil
}

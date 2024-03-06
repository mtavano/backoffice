package user

import (
	"time"

	"github.com/darchlabs/backoffice/internal/api/context"
	"github.com/darchlabs/backoffice/internal/storage"
	cardsdb "github.com/darchlabs/backoffice/internal/storage/cards"
	"github.com/darchlabs/backoffice/internal/storage/profile"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

type updateCardQuery func(storage.QueryContext, string) error

type PutProfileHandler struct {
	upsertProfileQuery  profileUpsertQuery
	selectCardByShortID selectCardsQuery
	updateCardQuery     updateCardQuery
}

type PutProfileRequest struct {
	UserID  string `json:"-"`
	ShortID string `json:"-"`

	Linkedin string `json:"linkedin"`
	Email    string `json:"email"`
	Whatsapp string `json:"whatsapp"`
	Medium   string `json:"medium"`
	TwitterX string `json:"twitterX"`
	Website  string `json:"website"`
}

func (h *PutProfileHandler) Invoke(ctx *context.Ctx, c *fiber.Ctx) (interface{}, int, error) {
	userID, err := context.GetUserIDFromRequestCtx(c)
	if err != nil {
		return nil, fiber.StatusUnauthorized, errors.New("Unauthorized")
	}

	shortID := c.Params("short_id")
	if shortID == "" {
		return nil, fiber.StatusBadRequest, errors.New("invalid operation. missing id")
	}

	var req PutProfileRequest
	err = c.BodyParser(&req)
	if err != nil {
		return nil, fiber.StatusBadRequest, errors.Wrap(err, "invalid body input")
	}

	req.UserID = userID
	req.ShortID = shortID

	return h.invoke(ctx, &req)
}

func (h *PutProfileHandler) invoke(
	ctx *context.Ctx, req *PutProfileRequest,
) (interface{}, int, error) {
	// query by short_id and return not allowed
	_, err := h.selectCardByShortID(ctx.App.SqlStore, req.ShortID)
	if errors.Is(err, cardsdb.ErrNoCard) {
		return nil, fiber.StatusForbidden, nil
	}
	if err != nil {
		return nil, fiber.StatusInternalServerError, nil
	}

	input := &profile.UpsertProfileInput{
		UserID:  req.UserID,
		ShortID: req.ShortID,
		Time:    time.Now(),
	}

	if req.Linkedin != "" {
		input.Linkedin = &req.Linkedin
	}
	if req.Email != "" {
		input.Email = &req.Email
	}
	if req.Whatsapp != "" {
		input.Whatsapp = &req.Whatsapp
	}
	if req.Medium != "" {
		input.Medium = &req.Medium
	}
	if req.Whatsapp != "" {
		input.Whatsapp = &req.Whatsapp
	}

	r, err := h.upsertProfileQuery(ctx.App.SqlStore, input)
	if err != nil {
		return nil, fiber.StatusInternalServerError, errors.Wrap(err, "cannot put profile record error")
	}

	err = h.updateCardQuery(ctx.App.SqlStore, req.ShortID)
	if err != nil {
		return nil, fiber.StatusInternalServerError, errors.Wrap(err, "cannot update card record error")
	}

	presented, status, err := ctx.PresentRecord(r, fiber.StatusOK)
	if err != nil {
		return nil, status, errors.Wrap(err, "upsert profile error")
	}

	return map[string]interface{}{
		"status":  "updated",
		"profile": presented,
	}, status, nil
}

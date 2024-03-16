package user

import (
	"log"
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

	Linkedin *string `json:"linkedin,omitempty" validate:"omitempty,nonempty"`
	Email    *string `json:"email,omitempty" validate:"omitempty,nonempty"`
	Whatsapp *string `json:"whatsapp,omitempty" validate:"omitempty,nonempty"`
	Medium   *string `json:"medium,omitempty" validate:"omitempty,nonempty"`
	TwitterX *string `json:"twitterX,omitempty" validate:"omitempty,nonempty"`
	Website  *string `json:"website,omitempty" validate:"omitempty,nonempty"`
	// Instagram *string `json:"instagram,omitempty" validate:"omitempty,nonempty"`
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
	err = ctx.Validator.Struct(req)
	if err != nil {
		return nil, fiber.StatusBadRequest, err
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
		log.Println("card not found")
		return nil, fiber.StatusForbidden, nil
	}
	if err != nil {
		log.Println("card error")
		return nil, fiber.StatusInternalServerError, nil
	}

	input, err := hydrateInput(req)
	if err != nil {
		return nil, fiber.StatusInternalServerError, errors.Wrap(err, "setting values error")
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

func hydrateInput(req *PutProfileRequest) (*profile.UpsertProfileInput, error) {
	input := &profile.UpsertProfileInput{
		UserID:      req.UserID,
		ShortID:     req.ShortID,
		Linkedin:    req.Linkedin,
		Email:       req.Email,
		Whatsapp:    req.Whatsapp,
		Medium:      req.Medium,
		TwitterX:    req.TwitterX,
		Website:     req.Website,
		Description: req.Description,
		Time:        time.Now(),
	}
	log.Printf("[handler] put profile input %+v", input)

	return input, nil
}

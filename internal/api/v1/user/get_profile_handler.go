package user

import (
	"database/sql"
	"time"

	"github.com/darchlabs/backoffice/internal/api/context"
	"github.com/darchlabs/backoffice/internal/storage/cards"
	"github.com/darchlabs/backoffice/internal/storage/profile"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

type GetProfileHandler struct {
	selectProfileQuery selectProfileQuery
	selectCardsQuery   selectCardsQuery
}

type getProfileHandlerResponse struct {
	Status   string `json:"status"`
	Nickname string `json:"nickname,omitemtpy"`
	ShortID  string `json:"shortId,omitemtpy"`

	// Social network links
	Linkedin *string `json:"linkedin,omitemtpy"`
	Email    *string `json:"email,omitemtpy"`
	Whatsapp *string `json:"whatsapp,omitemtpy"`
	Medium   *string `json:"medium,omitemtpy"`
	TwitterX *string `json:"twitterX,omitemtpy"`
	Website  *string `json:"website,omitemtpy"`

	// Non available fort the moment
	//Image string `json:"image"`

	CreatedAt time.Time  `json:"createdAt,omitemtpy"`
	UpdatedAt *time.Time `json:"updatedAt,omitemtpy"`
}

func (h *GetProfileHandler) Invoke(ctx *context.Ctx, c *fiber.Ctx) (interface{}, int, error) {
	shortID := c.Query("sid")
	nickname := c.Query("nn")

	card, err := h.selectCardsQuery(ctx.App.SqlStore, shortID)
	if errors.Is(err, cards.ErrNoCard) {
		return nil, fiber.StatusNotFound, nil
	}
	if err != nil {
		return nil, fiber.StatusInternalServerError, errors.Wrap(
			err, "user: GetProfileHandler.Invoke h.selectCardsQuery error",
		)
	}
	if card.Status == cards.StatusFree {
		return &getProfileHandlerResponse{
			Status: card.Status,
		}, fiber.StatusOK, nil
	}

	profile, err := h.selectProfileQuery(ctx.App.SqlStore, &profile.SelectFilters{
		ShortID:  shortID,
		Nickname: nickname,
	})
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fiber.StatusNotFound, nil
	}
	if err != nil {
		return nil, fiber.StatusInternalServerError, errors.Wrap(err, "something went wrong during operation error")
	}

	return &getProfileHandlerResponse{
		Status:    card.Status,
		Nickname:  profile.Nickname,
		ShortID:   profile.ShortID,
		Linkedin:  profile.Linkedin,
		Email:     profile.Email,
		Whatsapp:  profile.Whatsapp,
		Medium:    profile.Medium,
		TwitterX:  profile.TwitterX,
		Website:   profile.Website,
		CreatedAt: profile.CreatedAt,
	}, fiber.StatusOK, nil
}

package user

import (
	"fmt"
	"log"
	"time"

	"github.com/darchlabs/backoffice/internal/api/context"
	"github.com/darchlabs/backoffice/internal/storage/cards"
	cardsdb "github.com/darchlabs/backoffice/internal/storage/cards"
	profiledb "github.com/darchlabs/backoffice/internal/storage/profile"
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
	log.Printf("[pkg: user] GetProfileHandler.Invoke [short_id: '%s'] [nickname: '%s'] ", shortID, nickname)

	var profile *profiledb.Record
	var card *cardsdb.Record
	profile, err := h.selectProfileQuery(ctx.App.SqlStore, &profiledb.SelectFilters{
		ShortID:  shortID,
		Nickname: nickname,
	})
	if err != nil && !errors.Is(err, profiledb.ErrNoProfile) {
		fmt.Println("00000000000 ", err.Error())
		fmt.Println("00000000000 ", errors.Is(err, profiledb.ErrNoProfile))
		return nil, fiber.StatusInternalServerError, errors.Wrap(err, "something went wrong during operation error")
	}

	if profile == nil && shortID == "" {
		return nil, fiber.StatusNotFound, nil
	}

	var sid string
	switch true {
	case profile != nil:
		sid = profile.ShortID
		break
	case shortID != "":
		sid = shortID
		break
	}

	card, err = h.selectCardsQuery(ctx.App.SqlStore, sid)
	if errors.Is(err, cards.ErrNoCard) {
		return nil, fiber.StatusNotFound, nil
	}
	if err != nil {
		fmt.Println("HERE 2")
		return nil, fiber.StatusInternalServerError, errors.Wrap(
			err, "user: GetProfileHandler.Invoke h.selectCardsQuery error",
		)
	}
	if card.Status == cards.StatusFree {
		return &getProfileHandlerResponse{
			Status: card.Status,
		}, fiber.StatusOK, nil
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

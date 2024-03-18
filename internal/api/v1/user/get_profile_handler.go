package user

import (
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

type getProfileHandlerRequest struct {
	ShortID  string
	Nickname string
	UserID   string
}

type getProfileHandlerResponse struct {
	Status   string `json:"status"`
	ShortID  string `json:"shortId,omitemtpy"`
	Owner    bool   `json:"owner"`
	CanClaim bool   `json:"canClaim"` // TODO: fillme

	// Social network links
	Linkedin    *string `json:"linkedin,omitemtpy"`
	Email       *string `json:"email,omitemtpy"`
	Whatsapp    *string `json:"whatsapp,omitemtpy"`
	Medium      *string `json:"medium,omitemtpy"`
	TwitterX    *string `json:"twitterX,omitemtpy"`
	Website     *string `json:"website,omitemtpy"`
	Description *string `json:"description,omitempty"`
	Nickname    *string `json:"nickname,omitemtpy"`

	// Non available fort the moment
	//Image string `json:"image"`

	CreatedAt time.Time  `json:"createdAt,omitemtpy"`
	UpdatedAt *time.Time `json:"updatedAt,omitemtpy"`
}

func (h *GetProfileHandler) Invoke(ctx *context.Ctx, c *fiber.Ctx) (interface{}, int, error) {
	shortID := c.Query("sid")
	nickname := c.Query("nn")
	log.Printf("[pkg: user] GetProfileHandler.Invoke [short_id: '%s'] [nickname: '%s'] ", shortID, nickname)

	userID, _ := context.GetUserIDFromRequestCtx(c)

	return h.invoke(ctx, &getProfileHandlerRequest{
		ShortID:  shortID,
		Nickname: nickname,
		UserID:   userID,
	})

}

func (h *GetProfileHandler) invoke(ctx *context.Ctx, req *getProfileHandlerRequest) (interface{}, int, error) {
	var profile *profiledb.Record
	var card *cardsdb.Record
	profile, err := h.selectProfileQuery(ctx.App.SqlStore, &profiledb.SelectFilters{
		ShortID:  req.ShortID,
		Nickname: req.Nickname,
	})
	if err != nil && !errors.Is(err, profiledb.ErrNoProfile) {
		return nil, fiber.StatusInternalServerError, errors.Wrap(err, "something went wrong during operation error")
	}

	if profile == nil && req.ShortID == "" {
		return nil, fiber.StatusNotFound, nil
	}

	var sid string

	switch true {
	case profile != nil:
		sid = profile.ShortID
		break
	case req.ShortID != "":
		sid = req.ShortID
		break
	}

	card, err = h.selectCardsQuery(ctx.App.SqlStore, sid)
	if errors.Is(err, cards.ErrNoCard) {
		return nil, fiber.StatusNotFound, nil
	}
	if err != nil {
		return nil, fiber.StatusInternalServerError, errors.Wrap(
			err, "user: GetProfileHandler.Invoke h.selectCardsQuery error",
		)
	}
	if card.Status == cards.StatusFree {
		return &userProfileResponse{
			Status:  card.Status,
			ShortID: card.ShortID,
		}, fiber.StatusOK, nil
	}

	return &userProfileResponse{
		Status:    card.Status,
		CanClaim:  profile.UserID != req.UserID,
		Owner:     req.UserID == profile.UserID, // TODO: fix me
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

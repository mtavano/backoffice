package user

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/darchlabs/backoffice/internal/api/context"
	"github.com/darchlabs/backoffice/internal/storage/profile"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

type GetProfileHandler struct {
	selectProfileQuery selectProfileQuery
}

type getProfileHandlerResponse struct {
	Nickname string `json:"nickname"`
	ShortID  string `json:"shortId"`

	// Social network links
	Linkedin *string `json:"linkedin"`
	Email    *string `json:"email"`
	Whatsapp *string `json:"whatsapp"`
	Medium   *string `json:"medium"`
	TwitterX *string `json:"twitterX"`
	Website  *string `json:"website"`

	// Non available fort the moment
	//Image string `json:"image"`

	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
}

func (h *GetProfileHandler) Invoke(ctx *context.Ctx, c *fiber.Ctx) (interface{}, int, error) {
	shortID := c.Query("sid")
	nickname := c.Query("nn")
	fmt.Println(">>>>>> sid", shortID)
	fmt.Println(">>>>>> nn", nickname)

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
	fmt.Printf("1. %+v\n", profile)
	fmt.Printf("2. %+v\n", err)

	return &getProfileHandlerResponse{
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

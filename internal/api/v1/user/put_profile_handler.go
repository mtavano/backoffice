package user

import (
	"net/http"
	"time"

	"github.com/darchlabs/backoffice/internal/api/context"
	"github.com/darchlabs/backoffice/internal/storage/profile"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

type PutProfileHandler struct {
	profileUpsertQuery profileUpsertQuery
}

type PutProfileRequest struct {
	UserID  string `json:"-"`
	ShortID string `json:"short_id"`

	Linkedin *string `json:"linkedin"`
	Email    *string `json:"email"`
	Whatsapp *string `json:"whatsapp"`
	Medium   *string `json:"medium"`
	TwitterX *string `json:"twitterX"`
	Website  *string `json:"website"`
}

func (h *PutProfileHandler) Invoke(ctx *context.Ctx, c *fiber.Ctx) (interface{}, int, error) {
	userID, err := context.GetUserIDFromRequestCtx(c)
	if err != nil {
		return nil, fiber.StatusUnauthorized, errors.New("Unauthorized")
	}

	var req PutProfileRequest
	err = c.BodyParser(&req)
	if err != nil {
		return nil, fiber.StatusBadRequest, errors.Wrap(err, "bad request error")
	}

	if req.ShortID == "" {
		return nil, fiber.StatusBadRequest, errors.New("invalid operation. missing id")
	}

	req.UserID = userID

	return h.invoke(ctx, &req)
}

func (h *PutProfileHandler) invoke(
	ctx *context.Ctx, req *PutProfileRequest,
) (interface{}, int, error) {
	r, err := h.profileUpsertQuery(ctx.App.SqlStore, &profile.UpsertProfileInput{
		UserID:   req.UserID,
		ShortID:  req.ShortID,
		Linkedin: req.Linkedin,
		Email:    req.Email,
		Whatsapp: req.Whatsapp,
		Medium:   req.Medium,
		Website:  req.Website,
		Time:     time.Now(),
	})
	if err != nil {
		return nil, fiber.StatusInternalServerError, errors.Wrap(err, "cannot put record error")
	}

	return map[string]interface{}{
		"status":  "updated",
		"profile": r,
	}, http.StatusCreated, nil
}

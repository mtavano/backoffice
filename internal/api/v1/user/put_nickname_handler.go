package user

import (
	"github.com/darchlabs/backoffice/internal/api/context"
	"github.com/darchlabs/backoffice/internal/storage"
	"github.com/darchlabs/backoffice/internal/storage/profile"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

type updateNicknameQuery func(storage.Transaction, *profile.UpdateNicknameInput) (*profile.Record, error)

type PutNicknameHandler struct {
	updateNicknameQuery updateNicknameQuery
}

type PutNicknameRequest struct {
	UserID   string  `json:"-"`
	Nickname *string `json:"nickname,omitempty" validate:"omitempty,nonempty"`
	ShortID  string  `json:"-"`
}

func (h *PutNicknameHandler) Invoke(ctx *context.Ctx, c *fiber.Ctx) (interface{}, int, error) {
	userID, err := context.GetUserIDFromRequestCtx(c)
	if err != nil {
		return nil, fiber.StatusUnauthorized, errors.New("Unauthorized")
	}

	var req PutNicknameRequest
	err = c.BodyParser(&req)
	if err != nil {
		return nil, fiber.StatusBadRequest, errors.Wrap(err, "invalid body input")
	}
	req.ShortID = c.Params("short_id")
	req.UserID = userID

	err = ctx.Validator.Struct(req)
	if err != nil {
		return nil, fiber.StatusBadRequest, err
	}

	return h.invoke(ctx, &req)
}

func (h *PutNicknameHandler) invoke(ctx *context.Ctx, req *PutNicknameRequest) (interface{}, int, error) {
	r, err := h.updateNicknameQuery(ctx.App.SqlStore, &profile.UpdateNicknameInput{
		ShortID:  req.ShortID,
		UserID:   req.UserID,
		Nickname: req.Nickname,
	})
	if err != nil {
		return nil, fiber.StatusInternalServerError, errors.Wrap(err, "upsert profile error")
	}

	return map[string]interface{}{
		"status": "updated",
		"profile": &userProfileResponse{
			ShortID:     r.ShortID,
			Linkedin:    r.Linkedin,
			Email:       r.Email,
			Whatsapp:    r.Whatsapp,
			Medium:      r.Medium,
			TwitterX:    r.TwitterX,
			Website:     r.Website,
			Description: r.Description,
			Nickname:    r.Nickname,
		},
	}, fiber.StatusOK, nil
}

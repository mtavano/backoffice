package user

import (
	"github.com/darchlabs/backoffice/internal/api/context"
	profiledb "github.com/darchlabs/backoffice/internal/storage/profile"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

type GetNicknameAvailabilityHandler struct {
	selectProfileByNicknameQuery selectProfileByNicknameQuery
}

type getNicknameAvailabilityResponse struct {
	Available bool `json:"available"`
}

func (h *GetNicknameAvailabilityHandler) Invoke(ctx *context.Ctx, c *fiber.Ctx) (interface{}, int, error) {
	nn := c.Params("nickname")
	if nn == "" {
		return nil, fiber.StatusBadRequest, errors.New("invalid nickname")
	}

	_, err := h.selectProfileByNicknameQuery(ctx.App.SqlStore, nn)
	if errors.Is(err, profiledb.ErrNoProfile) {
		return &getNicknameAvailabilityResponse{
			Available: true,
		}, fiber.StatusOK, nil
	}
	if err != nil {
		return nil, fiber.StatusInternalServerError, errors.Wrap(
			err,
			"user: GetNicknameAvailabilityHandler.Invoke h.selectUserByNickname error",
		)
	}

	return &getNicknameAvailabilityResponse{
		Available: false,
	}, fiber.StatusOK, nil
}

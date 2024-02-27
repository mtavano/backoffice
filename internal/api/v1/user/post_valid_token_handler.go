package user

import (
	"encoding/json"

	"github.com/darchlabs/backoffice/internal/api/context"
	authdb "github.com/darchlabs/backoffice/internal/storage/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

type PostValidTokenHandler struct {
	secretKey              string
	authSelectByTokenQuery authSelectByTokenQuery
	userSelectByEmailQuery userSelectByEmailQuery
}

type PostValidTokenHandlerRequest struct {
	Token string `json:"token"`
}

type PostValidTokenHandlerResponse struct {
	UserID string `json:"user_id"`
}

func (h *PostValidTokenHandler) Invoke(ctx *context.Ctx, c *fiber.Ctx) (interface{}, int, error) {
	var req PostValidTokenHandlerRequest
	err := json.Unmarshal(c.Body(), &req)
	if err != nil {
		return nil, fiber.StatusBadRequest, errors.Wrap(err, "user: PostValidTokenHandler.Invoke c.BodyParser error")
	}

	payload, status, err := h.invoke(ctx, &req)
	if err != nil {
		return nil, status, errors.Wrap(err, "user: PostValidTokenHandler.Invoke h.invoke error")
	}
	return payload, status, nil
}

func (h *PostValidTokenHandler) invoke(ctx *context.Ctx, req *PostValidTokenHandlerRequest) (interface{}, int, error) {
	// Validate the token
	auth, err := h.authSelectByTokenQuery(ctx.App.SqlStore, req.Token)
	if errors.Is(err, authdb.ErrNotFound) {
		return nil, fiber.StatusUnauthorized, errors.Wrap(err, "user: PostValidTokenHandler.invoke h.authSelectByTokenQuery error")
	}
	if err != nil {
		return nil, fiber.StatusInternalServerError, errors.Wrap(err, "user: PostValidTokenHandler.invoke h.selectByTokenQuery error")
	}

	claims, err := parseToken(h.secretKey, auth.Token)
	if err != nil {
		return nil, fiber.StatusInternalServerError, errors.Wrap(err, "user: PostValidTokenHandler.invoke h.parseToken error")
	}

	userRecord, err := h.userSelectByEmailQuery(ctx.App.SqlStore, claims.Email)
	if err != nil {
		return nil, fiber.StatusInternalServerError, errors.Wrap(err, "user: PostValidTokenHandler.invoke h.userSelectByEmailQuery error")
	}

	// Perform additional checks on the claims, if necessary

	return &PostValidTokenHandlerResponse{UserID: userRecord.ID}, fiber.StatusOK, nil
}

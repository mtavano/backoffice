package user

import (
	"encoding/json"

	"github.com/darchlabs/backoffice/internal/api/context"
	authdb "github.com/darchlabs/backoffice/internal/storage/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

type PostValidApiKeyHandler struct {
	secretKey                string
	apikeySelectByTokenQuery apikeySelectByTokenQuery
}

type PostValidApiKeyHandlerRequest struct {
	ApiKey string `json:"api_key"`
}

type PostValidApiKeyHandlerResponse struct {
	UserID string `json:"user_id"`
}

func (h *PostValidApiKeyHandler) Invoke(ctx *context.Ctx, c *fiber.Ctx) (interface{}, int, error) {
	var req PostValidApiKeyHandlerRequest
	err := json.Unmarshal(c.Body(), &req)
	if err != nil {
		return nil, fiber.StatusBadRequest, errors.Wrap(err, "user: PostValidApiKeyHandler.Invoke c.BodyParser error")
	}

	payload, status, err := h.invoke(ctx, &req)
	if err != nil {
		return nil, status, errors.Wrap(err, "user: PostValidApiKeyHandler.Invoke h.invoke error")
	}
	return payload, status, nil
}

func (h *PostValidApiKeyHandler) invoke(ctx *context.Ctx, req *PostValidApiKeyHandlerRequest) (interface{}, int, error) {
	// Validate the token
	apiKey, err := h.apikeySelectByTokenQuery(ctx.App.SqlStore, req.ApiKey)
	if errors.Is(err, authdb.ErrNotFound) {
		return nil, fiber.StatusUnauthorized, errors.Wrap(err, "user: PostValidApiKeyHandler.invoke h.apikeySelectByTokenQuery error")
	}
	if err != nil {
		return nil, fiber.StatusInternalServerError, errors.Wrap(err, "user: PostValidApiKeyHandler.invoke h.apikeySelectByTokenQuery error")
	}

	claims, err := parseApiKey(h.secretKey, apiKey.Token)
	if err != nil {
		return nil, fiber.StatusInternalServerError, errors.Wrap(err, "user: PostValidApiKeyHandler.invoke h.parseToken error")
	}

	return &PostValidApiKeyHandlerResponse{UserID: claims.UserID}, fiber.StatusOK, nil
}

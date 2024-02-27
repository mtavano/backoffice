package user

import (
	"time"

	"github.com/darchlabs/backoffice/internal/api/context"
	"github.com/darchlabs/backoffice/internal/storage/apikey"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
)

type PostApiKeyHandler struct {
	secretKey              string
	idGenerate             idGenerate
	apikeyInsertQuery      apikeyInsertQuery
	authSelectByTokenQuery authSelectByTokenQuery
}

type postApiKeyHandlerRequest struct {
	UserID       string `json:"-"`
	DaysInterval int    `json:"days_interval"`
}

type postApiKeyHandlerResponse struct {
	ApiKey string `json:"api_key"`
}

type apiKeyClaims struct {
	UserID string `json:"email"`
	jwt.StandardClaims
}

func (h *PostApiKeyHandler) Invoke(ctx *context.Ctx, c *fiber.Ctx) (interface{}, int, error) {
	var req postApiKeyHandlerRequest
	err := c.BodyParser(&req)
	if err != nil {
		return nil, fiber.StatusInternalServerError, errors.Wrap(
			err, "api: PostApiKeyHandler.Invoke c.BodyParser error",
		)
	}

	id := c.Locals("user_id")
	userID, ok := id.(string)
	if !ok {
		return nil, fiber.StatusInternalServerError, errors.Wrap(err, "user: PostApiKeyHandler.Invoke id.(string) casting error")
	}
	req.UserID = userID

	payload, status, err := h.invoke(ctx, &req)
	if err != nil {
		return nil, status, errors.Wrap(err, "user: PostApiKeyHandler.Invoke h.invoke error")
	}

	return payload, status, nil
}

func (h *PostApiKeyHandler) invoke(ctx *context.Ctx, req *postApiKeyHandlerRequest) (interface{}, int, error) {
	// TODO: use better token valid interval
	now := time.Now()
	expiresAt := now.Add(24 * time.Duration(req.DaysInterval) * time.Hour)
	claims := apiKeyClaims{
		UserID: req.UserID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedApiKey, err := token.SignedString([]byte(h.secretKey))
	if err != nil {
		return nil, fiber.StatusInternalServerError, errors.Wrap(err, "auth: PostApiKeyHandler.invoke token.SignedString error")
	}

	err = h.apikeyInsertQuery(ctx.App.SqlStore, &apikey.Record{
		ID:        h.idGenerate(),
		UserID:    req.UserID,
		Token:     signedApiKey,
		ExpiresAt: &expiresAt,
		CreatedAt: &now,
	})
	if err != nil {
		return nil, fiber.StatusInternalServerError, errors.Wrap(err, "auth: PostApiKeyHandler.invoke h.apikeyInsertQuery error")
	}

	return &postApiKeyHandlerResponse{ApiKey: signedApiKey}, fiber.StatusCreated, nil
}

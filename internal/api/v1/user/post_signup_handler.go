package user

import (
	"time"

	"github.com/darchlabs/backoffice/internal/api/context"
	authdb "github.com/darchlabs/backoffice/internal/storage/auth"
	userdb "github.com/darchlabs/backoffice/internal/storage/user"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type PostSignupHandler struct {
	idGenerate      idGenerate
	secretKey       string
	userInsertQuery userInsertQuery
	authUpsertQuery authUpsertQuery
}

type postSignupHandlerRequest struct {
	Email    string `db:"email"`
	Password string `db:"password"`
}

type PostSignupHandlerResponse struct {
	Token string `json:"token"`
}

func (h *PostSignupHandler) Invoke(ctx *context.Ctx, c *fiber.Ctx) (interface{}, int, error) {
	var req postSignupHandlerRequest
	err := c.BodyParser(&req)
	if err != nil {
		return nil, fiber.StatusInternalServerError, errors.Wrap(
			err, "api: PostSignupHandler.Invoke c.BodyParser error",
		)
	}

	payload, status, err := h.invoke(ctx, &req)
	if err != nil {
		return nil, status, errors.Wrap(err, "user: PostSignupHandler.Invoke h.invoke error")
	}
	return payload, status, nil
}

func (h *PostSignupHandler) invoke(ctx *context.Ctx, req *postSignupHandlerRequest) (interface{}, int, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fiber.StatusInternalServerError, errors.Wrap(
			err, "user: PostSignupHandler.invoke bcrypt.GenerateFromPassword error",
		)
	}

	// Create *userdb.Record
	record := &userdb.Record{
		ID:             h.idGenerate(),
		Email:          req.Email,
		HashedPassword: string(hashedPassword),
		Verified:       false,
		CreatedAt:      time.Now(),
	}

	err = h.userInsertQuery(ctx.App.SqlStore, record)
	if err != nil {
		return nil, fiber.StatusInternalServerError, errors.Wrap(
			err, "user: PostSignupHandler.invoke h.userInsertQuery error",
		)
	}

	// TODO: use better token valid interval
	claims := loginClaims{
		Email: req.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * 365 * time.Hour).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(h.secretKey))
	if err != nil {
		return nil, fiber.StatusInternalServerError, errors.Wrap(err, "auth: PostLoginHandler.invoke token.SignedString error")
	}

	err = h.authUpsertQuery(ctx.App.SqlStore, &authdb.Record{
		UserID:    record.ID,
		Token:     signedToken,
		Blacklist: false,
		Kind:      authdb.TokenKindUser,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return nil, fiber.StatusInternalServerError, errors.Wrap(err, "auth: PostLoginHandler.invoke h.authUpsertQuery error")
	}

	return &PostSignupHandlerResponse{
		Token: signedToken,
	}, fiber.StatusCreated, nil
}

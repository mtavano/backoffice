package user

import (
	"time"

	"github.com/darchlabs/backoffice/internal/storage"
	"github.com/darchlabs/backoffice/internal/storage/apikey"
	"github.com/darchlabs/backoffice/internal/storage/auth"
	"github.com/darchlabs/backoffice/internal/storage/cards"
	"github.com/darchlabs/backoffice/internal/storage/profile"
	"github.com/darchlabs/backoffice/internal/storage/user"
	userdb "github.com/darchlabs/backoffice/internal/storage/user"
	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
)

type idGenerate func() string

type userInsertQuery func(storage.QueryContext, *userdb.Record) error

type userSelectByEmailQuery func(storage.Transaction, string) (*user.Record, error)

type authInsertQuery func(storage.QueryContext, *auth.Record) error

type authUpsertQuery func(storage.QueryContext, *auth.Record) error

type authSelectByTokenQuery func(storage.Transaction, string) (*auth.Record, error)

type apikeyInsertQuery func(storage.QueryContext, *apikey.Record) error

type apikeySelectByTokenQuery func(tx storage.Transaction, token string) (*apikey.Record, error)

type profileUpsertQuery func(storage.Transaction, *profile.UpsertProfileInput) (*profile.Record, error)

type selectProfileQuery func(storage.Transaction, *profile.SelectFilters) (*profile.Record, error)

type selectCardsQuery func(storage.Transaction, string) (*cards.Record, error)

func parseToken(secretKey, tokenString string) (*loginClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &loginClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, errors.Wrap(err, "jwt.ParseWithClaims error")
	}

	claims, ok := token.Claims.(*loginClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Verify that the token is not expired
	if claims.ExpiresAt < time.Now().Unix() {
		return nil, errors.New("token expired")
	}

	return claims, nil
}

func parseApiKey(secretKey, tokenString string) (*apiKeyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &apiKeyClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, errors.Wrap(err, "jwt.ParseWithClaims error")
	}

	claims, ok := token.Claims.(*apiKeyClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Verify that the token is not expired
	if claims.ExpiresAt < time.Now().Unix() {
		return nil, errors.New("token expired")
	}

	return claims, nil
}

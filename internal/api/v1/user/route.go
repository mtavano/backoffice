package user

import (
	"fmt"
	"net/http"

	"github.com/darchlabs/backoffice/internal/api/context"
	v1 "github.com/darchlabs/backoffice/internal/api/v1"
	authdb "github.com/darchlabs/backoffice/internal/storage/auth"
	cardsdb "github.com/darchlabs/backoffice/internal/storage/cards"
	profiledb "github.com/darchlabs/backoffice/internal/storage/profile"
	userdb "github.com/darchlabs/backoffice/internal/storage/user"
	"github.com/darchlabs/backoffice/pkg/client"
	"github.com/darchlabs/backoffice/pkg/middleware"
	"github.com/google/uuid"
)

func Route(basePath string, ctx *context.Ctx) {
	// handlers

	// SIGNUP
	postSignupHandler := &PostSignupHandler{
		secretKey:       ctx.App.Config.SecretKey,
		userInsertQuery: userdb.InsertQuery,
		idGenerate:      uuid.NewString,
		authUpsertQuery: authdb.UpsertQuery,
	}

	// LOGIN
	postLoginHandler := &PostLoginHandler{
		secretKey:              ctx.App.Config.SecretKey,
		userSelectByEmailQuery: userdb.SelectByEmailQuery,
		authUpsertQuery:        authdb.UpsertQuery,
	}

	postValidTokenHandler := &PostValidTokenHandler{
		secretKey:              ctx.App.Config.SecretKey,
		authSelectByTokenQuery: authdb.SelectByTokenQuery,
		userSelectByEmailQuery: userdb.SelectByEmailQuery,
	}

	putProfileHandler := &PutProfileHandler{
		upsertProfileQuery:  profiledb.UpsertQuery,
		selectCardByShortID: cardsdb.SelectQuery,
		updateCardQuery:     cardsdb.UpdateQuery,
	}

	getProfileHandler := &GetProfileHandler{
		selectProfileQuery: profiledb.SelectQuery,
		selectCardsQuery:   cardsdb.SelectQuery,
	}

	getNicknameAvailabilityHandler := &GetNicknameAvailabilityHandler{
		selectUserByNickname: userdb.SelectByNicknameQuery,
	}

	// setup middleware
	cl := client.New(&client.Config{
		Client:  http.DefaultClient,
		BaseURL: fmt.Sprintf("http://0.0.0.0:%s", ctx.App.Config.ApiServerPort),
	})

	auth := middleware.NewAuth(cl)

	// route
	ctx.Server.Get(
		fmt.Sprintf("%s/nn/:nickname", basePath),
		v1.HandleFunc(ctx, getNicknameAvailabilityHandler.Invoke),
	)
	ctx.Server.Post(
		fmt.Sprintf("%s/signup", basePath),
		v1.HandleFunc(ctx, postSignupHandler.Invoke),
	)
	ctx.Server.Post(
		fmt.Sprintf("%s/login", basePath),
		v1.HandleFunc(ctx, postLoginHandler.Invoke),
	)
	ctx.Server.Post(
		fmt.Sprintf("%s/tokens", basePath),
		v1.HandleFunc(ctx, postValidTokenHandler.Invoke),
	)
	fmt.Printf("%s/profiles\n", basePath)
	ctx.Server.Get(
		fmt.Sprintf("%s/profiles", basePath),
		v1.HandleFunc(ctx, getProfileHandler.Invoke),
	)
	ctx.Server.Put(
		fmt.Sprintf("%s/profiles/:short_id", basePath),
		auth.Middleware,
		v1.HandleFunc(ctx, putProfileHandler.Invoke),
	)
}

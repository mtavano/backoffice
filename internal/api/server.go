package api

import (
	"fmt"
	"os"

	"github.com/darchlabs/backoffice/internal/api/admin"
	"github.com/darchlabs/backoffice/internal/api/context"
	v1 "github.com/darchlabs/backoffice/internal/api/v1"
	"github.com/darchlabs/backoffice/internal/api/v1/user"
	"github.com/darchlabs/backoffice/internal/application"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type ServerConfig struct {
	Port string
	App  *application.App
}

type Server struct {
	server *fiber.App
	app    *application.App
	port   string
}

func NewServer(config *ServerConfig) *Server {
	server := fiber.New()
	server.Use(logger.New())
	server.Use(logger.New(logger.Config{
		Format:     "[${ip}]:${port} ${status} - ${method} ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
		Output:     os.Stdout,
	}))
	// Or extend your config for customization
	server.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "*",
	}))

	return &Server{
		server: server,
		app:    config.App,
		port:   config.Port,
	}
}

func (s *Server) Start(app *application.App) error {
	go func() {
		ctx := context.New(&context.Config{
			Server: s.server,
			App:    s.app,
		})
		// route endpoints
		v1.HealthRoute(ctx)
		user.Route("/api/v1/users", ctx)
		admin.Route("/admin", ctx)

		s.server.Get("/api/v1/health", v1.HandleFunc(
			ctx,
			func(_ *context.Ctx, _ *fiber.Ctx) (interface{}, int, error) {
				return map[string]string{"status": "running"}, fiber.StatusOK, nil
			},
		))

		// sever listen
		err := s.server.Listen(fmt.Sprintf(":%s", s.port))
		if err != nil {
			panic(err)
		}
	}()

	return nil
}

package context

import (
	"github.com/darchlabs/backoffice/internal/application"
	"github.com/gofiber/fiber/v2"
)

type Ctx struct {
	// structs
	Server *fiber.App
	App    *application.App

	// interfaces
	//SqlStore storage.SQL
}

type Config struct {
	// stucts
	Server *fiber.App
	App    *application.App

	// interfaces
	//SqlStore storage.SQL
}

func New(conf *Config) *Ctx {
	return &Ctx{
		Server: conf.Server,
		App:    conf.App,
		//SqlStore: conf.SqlStore,
	}
}

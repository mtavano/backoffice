package context

import (
	"github.com/darchlabs/backoffice/internal/application"
	"github.com/gofiber/fiber/v2"
	"github.com/teris-io/shortid"
)

type shortIDGen func() (string, error)

type Ctx struct {
	// structs
	Server           *fiber.App
	App              *application.App
	ShortIDGenerator shortIDGen
}

type Config struct {
	// stucts
	Server *fiber.App
	App    *application.App
}

func New(conf *Config) *Ctx {
	return &Ctx{
		Server:           conf.Server,
		App:              conf.App,
		ShortIDGenerator: shortid.Generate,
	}
}

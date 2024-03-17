package context

import (
	"github.com/darchlabs/backoffice/internal/application"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/teris-io/shortid"
)

type shortIDGen func() (string, error)

type Ctx struct {
	// structs
	Server           *fiber.App
	App              *application.App
	ShortIDGenerator shortIDGen
	Validator        *validator.Validate
}

type Config struct {
	// stucts
	Server *fiber.App
	App    *application.App
}

func New(conf *Config) (*Ctx, error) {
	vv := validator.New()
	vv.RegisterValidation("nonempty", nonEmpty)

	return &Ctx{
		Server:           conf.Server,
		App:              conf.App,
		ShortIDGenerator: shortid.Generate,
		Validator:        vv,
	}, nil
}

// custom validatin func
func nonEmpty(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	return len(val) > 0
}

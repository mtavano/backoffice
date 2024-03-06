package context

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/darchlabs/backoffice/internal/application"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
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

func (c *Ctx) PresentRecord(v interface{}, status int) (interface{}, int, error) {
	bb, err := json.Marshal(v)
	if err != nil {
		return nil, fiber.StatusInternalServerError, errors.Wrap(err, "cannot present data properly error")
	}

	var result map[string]interface{}
	if err := json.Unmarshal(bb, &result); err != nil {
		return nil, fiber.StatusInternalServerError, errors.Wrap(err, "cannot re-present data properly error")
	}

	transformed := make(map[string]interface{})
	for key, value := range result {
		k := fmt.Sprintf("%s%s", strings.ToLower(key[:1]), key[1:])
		transformed[k] = value
	}

	return transformed, status, nil
}

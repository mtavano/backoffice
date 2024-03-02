package middleware

import "github.com/gofiber/fiber/v2"

type AdminKey struct {
	adminKey string
}

func NewAdminKey(key string) *AdminKey {
	return &AdminKey{adminKey: key}
}

func (a *AdminKey) Middleware(c *fiber.Ctx) error {
	adminKey := c.Get("x-admin-key")
	if adminKey != a.adminKey {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}

	return c.Next()
}

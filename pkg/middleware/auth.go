package middleware

import (
	"strings"

	"github.com/darchlabs/backoffice/pkg/client"
	"github.com/gofiber/fiber/v2"
)

type Auth struct {
	client *client.Client
}

func NewAuth(cl *client.Client) *Auth {
	return &Auth{client: cl}
}

func (a *Auth) Middleware(strict bool) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// Extract the Authorization header
		authHeader := c.Get("Authorization")

		// Check if the Authorization header is present and in the correct format
		//if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		//return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		//}

		// Extract the token from the Authorization header
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Call the ValidTokenWithCtx function to validate the token and get the user ID
		response, err := a.client.ValidTokenWithCtx(c.Context(), token)
		if err != nil && strict {
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}

		// Attach the user ID to the request context
		userID := ""
		if response != nil {
			userID = response.UserID
		}

		c.Locals("user_id", userID)

		// Continue to the next middleware or route handler
		return c.Next()
	}
}

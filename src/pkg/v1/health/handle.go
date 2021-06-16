package health

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// RegisterHealth handles the registration of all health related routes
func RegisterHealth(r fiber.Router) {
	r.Get("/health", func(ctx *fiber.Ctx) error {
		return ctx.JSON(map[string]interface{}{
			"status":  http.StatusOK,
			"message": "OK",
		})
	})
}

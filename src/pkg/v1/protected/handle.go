package protected

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// RegisterProtected will serve as a test endpoint for jwt middlewares
func RegisterProtected(r fiber.Router) {
	protectedGroup := r.Group("/protected")
	protectedGroup.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.JSON(map[string]interface{}{
			"status":  http.StatusOK,
			"message": "OK",
		})
	})
}

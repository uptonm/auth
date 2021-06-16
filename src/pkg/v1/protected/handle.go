package protected

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func RegisterProtected(r fiber.Router) {
	protectedGroup := r.Group("/protected")
	protectedGroup.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.JSON(map[string]interface{}{
			"status":  http.StatusOK,
			"message": "OK",
		})
	})
}

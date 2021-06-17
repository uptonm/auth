package protected

import (
	"github.com/uptonm/auth/src/common"
	"github.com/uptonm/auth/src/pkg/auth"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func RegisterProtected(r fiber.Router, config *common.Config) {
	protectedGroup := r.Group("/protected")
	protectedGroup.Use(auth.CreateJwtMiddleware(config))
	protectedGroup.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.JSON(map[string]interface{}{
			"status":  http.StatusOK,
			"message": "OK",
		})
	})
}

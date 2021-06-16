package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/shareed2k/goth_fiber"
	log "github.com/sirupsen/logrus"
)

func RegisterAuthRoutes(r fiber.Router) {
	authGroup := r.Group("/auth")

	authGroup.Get("/:provider", goth_fiber.BeginAuthHandler)
	authGroup.Get("/:provider/callback", func(ctx *fiber.Ctx) error {
		user, err := goth_fiber.CompleteUserAuth(ctx)
		if err != nil {
			return fiber.ErrUnauthorized
		}

		return ctx.Render("authenticated", user)
	})

	authGroup.Get("/:provider/logout", func(ctx *fiber.Ctx) error {
		if err := goth_fiber.Logout(ctx); err != nil {
			log.Fatal(err)
		}
		return ctx.Redirect("/")
	})
}

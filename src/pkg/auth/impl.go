package auth

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"

	"github.com/uptonm/auth/src/common"
)

func CreateJwtMiddleware(config *common.Config) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:    []byte(config.SigningKey),
		SigningMethod: "HS256",
	})
}

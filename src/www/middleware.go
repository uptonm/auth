package www

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

// wireMiddleware handles the initialization of the middleware chain
func wireMiddleware(r fiber.Router) {
	r.Use(recover.New())

	r.Use(requestid.New(requestid.Config{
		Header:     "X-Request-ID",
		ContextKey: "x_request_id",
	}))

	r.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:8080, https://localhost:8080",
		AllowMethods: "GET",
	}))

	r.Use(logger.New())
}

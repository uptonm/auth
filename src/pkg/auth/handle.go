package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"

	"github.com/uptonm/auth/src/pkg/db"
)

// RegisterAuthRoutes handles registration of all authentication related handlers
func RegisterAuthRoutes(r fiber.Router) {
	authGroup := r.Group("/auth")
	authGroup.Get("/", HandleSignIn)
	authGroup.Get("/callback", HandleSignInCallback)
}

// HandleSignIn redirects the user to auth0 with the proper query parameters
// to begin an authorization request
func HandleSignIn(ctx *fiber.Ctx) error {
	state, err := generateStateFromCtx(ctx)
	if err != nil {
		log.Errorf("failed to generate state from fiber context")
		return fiber.ErrInternalServerError
	}

	err = db.RedisConn.Set(context.Background(), fmt.Sprintf("state:%s", state), "true", time.Minute*5).Err()
	if err != nil {
		log.Errorf("failed to store state key in redis")
		return fiber.ErrInternalServerError
	}

	return ctx.Redirect(buildAuthorizationUrl(authRequestConfig{
		State:        state,
		Scope:        "openid profile",
		Connection:   nil,
		Organization: nil,
		Invitation:   nil,
	}), fiber.StatusTemporaryRedirect)
}

// HandleSignInCallback handles the auth0 redirect, validating the state parameter
// to prevent CSRF and using the callback code supplied to authenticate the user,
// generate access/refresh tokens, and return them as the result
func HandleSignInCallback(ctx *fiber.Ctx) error {
	code := ctx.Query("code")
	if len(code) == 0 {
		log.Errorf("authentication callback contained no code parameter")
		return fiber.ErrBadRequest
	}

	state := ctx.Query("state")
	if len(state) == 0 {
		log.Errorf("authentication callback contained no state parameter")
		return fiber.ErrBadRequest
	}

	res, err := db.RedisConn.Get(context.Background(), fmt.Sprintf("state:%s", state)).Result()
	if err != nil {
		log.Errorf("failed to query redis for state parameter")
		return fiber.ErrInternalServerError
	}

	if res != "true" {
		log.Errorf("failed to validate request state parameter")
		return fiber.ErrUnauthorized
	}

	payload, err := generateRequestToken(code)
	if err != nil {
		log.Errorf("failed to generate request token")
		return fiber.ErrUnauthorized
	}

	go invalidateStateSync(state)

	return ctx.JSON(payload)
}

// HandleSignOut will eventually invalidate a jwt and refresh token
func HandleSignOut(ctx *fiber.Ctx) error {
	return ctx.Redirect("/")
}

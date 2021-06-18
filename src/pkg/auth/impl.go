package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"github.com/uptonm/auth/src/pkg/db"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/uptonm/auth/src/common"
)

const auth0BaseUrl = "https://%s/authorize"

type authRequestConfig struct {
	// State is an opaque arbitrary alphanumeric string your app adds to the initial
	// request that Auth0 includes when redirecting back to your application. This is
	// verified on the callback to prevent cross-site request forgery (CSRF) attacks.
	State string `json:"state"`
	// Scope Specifies the scopes for which you want to request authorization,
	// which dictate which claims (or user attributes) you want returned.
	Scope string `json:"scope"`
	// Connection forces the user to sign in with a specific connection. For example,
	// you can pass a value of github to send the user directly to GitHub to log in with
	// their GitHub account. When not specified, the user sees the Auth0 Lock screen with
	// all configured connections.
	Connection *string `json:"connection"`
	// Organization represents the ID of the organization to use when authenticating a user.
	// When not provided, if your application is configured to Display Organization Prompt,
	// the user will be able to enter the organization name when authenticating.
	Organization *string `json:"organization"`
	// Invitation is the Ticket ID of the organization invitation. When inviting a member
	// to an Organization, your application should handle invitation acceptance by forwarding
	// the invitation and organization key-value pairs when the user accepts the invitation.
	Invitation *string `json:"invitation"`
}

// buildAuthorizationUrl utilizes a supplied authRequestConfig to generate an
// auth0 authorization code baseUrl to start implicit flow authentication
func buildAuthorizationUrl(requestConfig authRequestConfig) string {
	authUrl := fmt.Sprintf(auth0BaseUrl, common.Config.Auth0Domain)
	params := url.Values{
		"response_type": []string{"code"},
		"state":         []string{requestConfig.State},
		"scope":         []string{requestConfig.Scope},
	}

	if requestConfig.Organization != nil {
		params.Set("organization", *requestConfig.Organization)
	}

	if requestConfig.Connection != nil {
		params.Set("connection", *requestConfig.Connection)
	}

	if requestConfig.Invitation != nil {
		params.Set("invitation", *requestConfig.Invitation)
	}

	return fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&%s", authUrl, common.Config.Auth0ClientId,
		common.Config.Auth0CallbackUrl, params.Encode())
}

// getClientIp uses a fiber.Ctx to read the X-Forwarded-For / X-Real-Ip headers
// to determine the client's IP
func getClientIp(ctx *fiber.Ctx) string {
	clientIP := ctx.Get("X-Forwarded-For")
	if index := strings.IndexByte(clientIP, ','); index >= 0 {
		clientIP = clientIP[0:index]
		// Get the first one, ie 1.1.1.1
	}
	clientIP = strings.TrimSpace(clientIP)
	if len(clientIP) > 0 {
		return clientIP
	}
	clientIP = strings.TrimSpace(string(ctx.Get("X-Real-Ip")))
	if len(clientIP) > 0 {
		return clientIP
	}
	return ctx.Context().RemoteIP().String()
}

// generateStateFromCtx uses the fiber.Ctx to retrieve the user's IP
// and base64 encode it for use as state parameter
func generateStateFromCtx(ctx *fiber.Ctx) (string, error) {
	data := []byte(getClientIp(ctx))
	if _, err := io.ReadFull(rand.Reader, data); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

type requestTokenPayload struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IdToken      string `json:"id_token"`
	TokenType    string `json:"token_type"`
}

// generateRequestToken accepts a callback code and performs the second step of auth code flow,
// sending a POST request to auth0 to generate an access token and refresh token
func generateRequestToken(code string) (*requestTokenPayload, error) {
	var payload requestTokenPayload
	requestUrl := fmt.Sprintf("https://%s/oauth/token", common.Config.Auth0Domain)

	requestPayload := strings.NewReader(
		fmt.Sprintf("grant_type=authorization_code&client_id=%s&client_secret=%s&code=%s&redirect_uri=%s",
			common.Config.Auth0ClientId, common.Config.Auth0ClientSecret, code, common.Config.Auth0CallbackUrl),
	)

	req, _ := http.NewRequest("POST", requestUrl, requestPayload)

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Errorf("failed to close io.Reader")
		}
	}(res.Body)
	body, _ := ioutil.ReadAll(res.Body)

	err = json.Unmarshal(body, &payload)
	if err != nil {
		log.Errorf("failed to unmarshal response body")
		return nil, err
	}

	return &payload, nil
}

// invalidateStateSync is a function used to remove a state key from redis
// without returning anything so it can be run in a go routine without blocking
// the response of the authentication request.
func invalidateStateSync(state string) {
	err := db.RedisConn.Del(context.Background(), fmt.Sprintf("state:%s", state)).Err()
	if err != nil {
		log.Errorf("unable to evict state key error=%s", err.Error())
	}
}

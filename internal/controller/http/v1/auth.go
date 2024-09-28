package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/magmaheat/auth_service/internal/service"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

type AuthRouter struct {
	Auth service.Auth
}

type UpdateRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type Response struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewAuthRouter(services service.Auth) *AuthRouter {
	return &AuthRouter{Auth: services}
}

func (a *AuthRouter) LoginHandler(c echo.Context) error {
	params := c.QueryParams()
	userId := params.Get("user_id")
	if userId == "" {
		NewErrorResponce(c, http.StatusBadRequest, ErrInvalidAuthHeader.Error())
		return ErrInvalidAuthHeader
	}

	userIp := c.RealIP()

	accessToken, refreshToken, err := a.Auth.GenerateTokens(userId, userIp)
	if err != nil {
		NewErrorResponce(c, http.StatusInternalServerError, "error generate tokens")
		return ErrGenerateTokens
	}

	return c.JSON(http.StatusCreated, Response{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
	return nil
}

func (a *AuthRouter) UpdateHandler(c echo.Context) error {
	token, ok := bearerToken(c.Request())
	if !ok {
		log.Errorf("http - v1 - auth - UpdateHandler: bearerToken: %v", ErrInvalidAuthHeader)
		NewErrorResponce(c, http.StatusUnauthorized, ErrInvalidAuthHeader.Error())
		return ErrInvalidAuthHeader
	}

	var input UpdateRequest
	if err := c.Bind(&input); err != nil {
		NewErrorResponce(c, http.StatusBadRequest, err.Error())
		return err
	}

	userIp := c.RealIP()

	accessToken, refreshToken, err := a.Auth.UpdateTokens(token, input.RefreshToken, userIp)
	if err != nil {
		NewErrorResponce(c, http.StatusUnauthorized, err.Error())
	}

	c.JSON(http.StatusOK, Response{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})

	return nil
}

func bearerToken(r *http.Request) (string, bool) {
	const prefix = "Bearer "

	header := r.Header.Get(echo.HeaderAuthorization)
	if header == "" {
		return "", false
	}

	if len(header) > len(prefix) && strings.EqualFold(header[:len(prefix)], prefix) {
		return header[len(prefix):], true
	}

	return "", false
}

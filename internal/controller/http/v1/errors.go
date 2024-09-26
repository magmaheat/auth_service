package v1

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
)

var (
	ErrInvalidAuthHeader = fmt.Errorf("invalid auth header")
	ErrCannotParseToken  = fmt.Errorf("cannot parse token")
)

func NewErrorResponce(c echo.Context, errStatus int, message string) {
	err := errors.New(message)
	httpErr, ok := err.(*echo.HTTPError)
	if !ok {
		httpErr = echo.NewHTTPError(errStatus, err.Error())
	}

	c.JSON(errStatus, httpErr)
}

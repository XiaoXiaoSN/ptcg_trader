package echo

import (
	"net/http"

	"ptcg_trader/internal/errors"

	"github.com/labstack/echo/v4"
)

// notFoundHandler responds not found response.
func notFoundHandler(c echo.Context) error {
	return c.JSON(http.StatusNotFound, errors.GetHTTPError(errors.ErrPageNotFound))
}

// httpErrorHandler responds error response according to given error form echo service.
func httpErrorHandler(err error, c echo.Context) {
	if err == nil {
		return
	}

	// handle echo HTTPError
	echoErr, ok := err.(*echo.HTTPError)
	if ok {
		_ = c.JSON(echoErr.Code, echoErr)
		return
	}

	// handle platform custom error
	causeErr := errors.Cause(err)
	httpErr := errors.GetHTTPError(causeErr)

	_ = c.JSON(httpErr.Status, httpErr)
}

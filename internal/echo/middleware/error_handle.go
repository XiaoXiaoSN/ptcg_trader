package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// ErrMiddleware returns the error handle middleware
func ErrMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)

		if err == nil {
			return nil
		}

		// record the error
		logFields := map[string]interface{}{}
		req := c.Request()
		{
			logFields["request_method"] = req.Method
			logFields["request_url"] = req.URL.String()
		}

		// dump response body
		resp := c.Response()
		resp.After(func() {
			logFields["response_status"] = resp.Status
			logger := log.Ctx(req.Context()).With().Fields(logFields).Logger()

			// use response status code to set log level
			if resp.Status >= http.StatusInternalServerError {
				logger.Error().Msgf("%+v", err)
			} else if resp.Status >= http.StatusBadRequest {
				logger.Warn().Msgf("%+v", err)
			} else {
				logger.Debug().Msgf("%+v", err)
			}
		})

		return err
	}

}

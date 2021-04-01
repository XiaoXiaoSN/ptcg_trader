package middleware

import (
	"context"
	"strconv"

	"ptcg_trader/internal/ctxutil"
	"ptcg_trader/internal/errors"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// IdentityMiddleware returns the identity middleware
func IdentityMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			// try to get identity from http header.
			identity := c.Request().Header.Get(ctxutil.CtxKeyIdentityID.String())
			if identity == "" {
				return errors.Wrap(errors.ErrUnauthorized, "CtxKeyIdentityID not exist in http request header")
			}
			identityID, err := strconv.ParseInt(identity, 10, 64)
			if err != nil {
				return errors.WithMessage(errors.ErrUnauthorized, "identity id not a number")
			}

			// context.Context with the identity
			ctx := context.WithValue(c.Request().Context(), ctxutil.CtxKeyIdentityID, identityID)

			// set the identity ctx into log system
			logger := log.Ctx(ctx).With().Int64("identity_id", identityID).Logger()
			ctx = logger.WithContext(ctx)

			// write back to echo request object (with identity ctx and logger)
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}

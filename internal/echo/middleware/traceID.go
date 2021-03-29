package middleware

import (
	"context"

	"ptcg_trader/internal/ctxutil"

	"github.com/labstack/echo/v4"
	"github.com/rs/xid"
	"github.com/rs/zerolog/log"
)

// NewTraceIDMiddleware returns the traceID middleware
func NewTraceIDMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			// try to get traceID from http header. or not, create one immediately
			traceID := c.Request().Header.Get(ctxutil.CtxKeyTraceID.String())
			if traceID == "" {
				traceID = xid.New().String()
			}

			// context.Context with the traceID
			ctx := context.WithValue(c.Request().Context(), ctxutil.CtxKeyTraceID, traceID)

			// set the traceID ctx into log system
			logger := log.With().Str("trace_id", traceID).Logger()
			ctx = logger.WithContext(ctx)

			// write back to echo request object (with traceID ctx and logger)
			c.SetRequest(c.Request().WithContext(ctx))

			// write traceID into echo response object
			c.Response().Writer.Header().Set(ctxutil.CtxKeyTraceID.String(), traceID)

			return next(c)
		}
	}
}

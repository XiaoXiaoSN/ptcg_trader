package middleware

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/rs/xid"
	"github.com/rs/zerolog/log"
)

// CtxKey define the key for context.Context
type CtxKey string

func (ck CtxKey) String() string {
	return string(ck)
}

// CtxKeyTraceID define http header name for TraceID
const CtxKeyTraceID CtxKey = "X-Trace-Id"

// TraceIDFromContext get the traceID form context.
// if that not exist, create one immediately
func TraceIDFromContext(ctx context.Context) string {
	rid, ok := ctx.Value(CtxKeyTraceID).(string)
	if !ok {
		rid = xid.New().String()
		return rid
	}
	return rid
}

// NewTraceIDMiddleware returns the traceID middleware
func NewTraceIDMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			// try to get traceID from http header. or not, create one immediately
			traceID := c.Request().Header.Get(CtxKeyTraceID.String())
			if traceID == "" {
				traceID = xid.New().String()
			}

			// context.Context with the traceID
			ctx := context.WithValue(c.Request().Context(), CtxKeyTraceID, traceID)

			// set the traceID ctx into log system
			logger := log.With().Str("trace_id", traceID).Logger()
			ctx = logger.WithContext(ctx)

			// write back to echo request object (with traceID ctx and logger)
			c.SetRequest(c.Request().WithContext(ctx))

			// write traceID into echo response object
			c.Response().Writer.Header().Set(CtxKeyTraceID.String(), traceID)

			return next(c)
		}
	}
}

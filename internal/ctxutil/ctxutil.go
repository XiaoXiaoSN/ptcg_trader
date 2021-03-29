package ctxutil

import (
	"context"
	"ptcg_trader/internal/errors"
	"strconv"

	"github.com/rs/xid"
)

// CtxKey define the key for context.Context
type CtxKey string

func (ck CtxKey) String() string {
	return string(ck)
}

// define some CtxKeys
const (
	// CtxKeyTraceID define http header name for TraceID
	CtxKeyTraceID CtxKey = "X-Trace-Id"

	// CtxKeyIdentityID define http header name for identityID, known as user id
	CtxKeyIdentityID CtxKey = "X-Identity-Id"
)

// TraceIDFromCtx get the traceID form context.
// if that not exist, create one immediately
func TraceIDFromCtx(ctx context.Context) string {
	rid, ok := ctx.Value(CtxKeyTraceID).(string)
	if !ok {
		rid = xid.New().String()
		return rid
	}
	return rid
}

// IdentityIDFromCtx get the identity form context.
func IdentityIDFromCtx(ctx context.Context) (int64, error) {
	switch identityID := ctx.Value(CtxKeyIdentityID).(type) {
	case int:
		return int64(identityID), nil
	case int64:
		return int64(identityID), nil
	case float64:
		return int64(identityID), nil
	case string:
		id, err := strconv.ParseInt(identityID, 10, 64)
		if err != nil {
			return 0, errors.Wrap(errors.ErrUnauthorized, "identity string not a number")
		}
		return id, nil
	}
	return 0, errors.ErrUnauthorized
}

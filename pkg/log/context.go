package log

import (
	"context"

	"github.com/rs/zerolog"
)

// Extendable list of generic, well-known and wildly used context keys.
// A particular logger could be automatically contextualized by these keys via Contextualize() method.
type (
	RequestIDContextKey   struct{}
	TraceIDContextKey     struct{}
	AuthSubjectContextKey struct{}
)

// AddContext adds context to a particular log event by predefined context keys.
func AddContext(ctx context.Context, event *zerolog.Event) *zerolog.Event {
	if v := ctx.Value(TraceIDContextKey{}); v != nil {
		event.Str("trace_id", strOrInvalid(v))
	}

	if v := ctx.Value(RequestIDContextKey{}); v != nil {
		event.Str("request_id", strOrInvalid(v))
	}

	if v := ctx.Value(AuthSubjectContextKey{}); v != nil {
		event.Str("authsub", strOrInvalid(v))
	}

	return event
}

func strOrInvalid(v any) string {
	x, ok := v.(string)
	if !ok {
		return "invalid"
	}

	return x
}

package transport

import "context"

type Header interface {
	Get(key string) string
	Set(key string, value string)
	Keys() []string
}

type Transporter interface {
	Operation() string
	RequestHeader() Header
}

type (
	serverTransportKey struct{}
)

// NewServerContext returns a new Context that carries value.
func NewServerContext(ctx context.Context, tr Transporter) context.Context {
	return context.WithValue(ctx, serverTransportKey{}, tr)
}

// FromServerContext returns the Transport value stored in ctx, if any.
func FromServerContext(ctx context.Context) (tr Transporter, ok bool) {
	tr, ok = ctx.Value(serverTransportKey{}).(Transporter)
	return
}

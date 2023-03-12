package servertiming

import (
	"context"

	orig "github.com/mitchellh/go-server-timing"
)

func NewContext(ctx context.Context, h *Header) context.Context {
	return orig.NewContext(ctx, h.Header)
}

func FromContext(ctx context.Context) *Header {
	h := orig.FromContext(ctx)
	if h == nil {
		return nil
	}

	return &Header{h}
}

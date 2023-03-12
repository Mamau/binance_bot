package servertiming

import (
	"net/http"

	orig "github.com/mitchellh/go-server-timing"
)

const HeaderKey = orig.HeaderKey

type Header struct {
	*orig.Header
}

func (h *Header) NewMetric(name string) *Metric {
	return h.Add(&Metric{&orig.Metric{Name: name}})
}

func (h *Header) Add(m *Metric) *Metric {
	return &Metric{h.Header.Add(m.Metric)}
}

func ParseHeader(input string) (*Header, error) {
	h, err := orig.ParseHeader(input)
	return &Header{h}, err
}

func writeHeader(headers http.Header, h *Header) {
	h.Lock()
	defer h.Unlock()

	if len(h.Metrics) == 0 {
		return
	}

	headers.Set(HeaderKey, h.String())
}

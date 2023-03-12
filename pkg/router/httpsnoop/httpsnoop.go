package httpsnoop

import (
	"github.com/gin-gonic/gin"
)

type WriteHeaderFunc func(code int)
type WriteFunc func(b []byte) (int, error)

type Hooks struct {
	Write       func(WriteFunc) WriteFunc
	WriteHeader func(WriteHeaderFunc) WriteHeaderFunc
}

type snoopingWriter struct {
	gin.ResponseWriter
	hooks Hooks
}

func Wrap(w gin.ResponseWriter, hooks Hooks) *snoopingWriter {
	return &snoopingWriter{
		ResponseWriter: w,
		hooks:          hooks,
	}
}

func (sw *snoopingWriter) Write(b []byte) (int, error) {
	if sw.hooks.Write != nil {
		return sw.hooks.Write(sw.ResponseWriter.Write)(b)
	}

	return sw.ResponseWriter.Write(b)
}

func (sw *snoopingWriter) WriteHeader(status int) {
	if sw.hooks.WriteHeader != nil {
		sw.hooks.WriteHeader(sw.ResponseWriter.WriteHeader)(status)
	}

	sw.ResponseWriter.WriteHeader(status)
}

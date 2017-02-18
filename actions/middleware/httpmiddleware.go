package middleware

import (
	"github.com/gobuffalo/buffalo"
	"net/http"
	"os"
)

type HttpHandlerWrapper struct {
	h     buffalo.Handler
	c     buffalo.Context
	error error
}

func (mw HttpHandlerWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	os.Stderr.WriteString("wrapper::before")
	mw.error = mw.h(mw.c)
	os.Stderr.WriteString("wrapper::after")
}

func HttpMiddleware(handler func(http.Handler) http.Handler) func(buffalo.Handler) buffalo.Handler {
	return func(next buffalo.Handler) buffalo.Handler {

		wrapper := HttpHandlerWrapper{h: next}

		return func(c buffalo.Context) error {
			wrapper.c = c
			current := handler(wrapper)
			current.ServeHTTP(c.Response(), c.Request())
			return wrapper.error
		}
	}
}

package middleware

import (
	"github.com/gobuffalo/buffalo"
	"net/http"
)

type httpMiddlewareFunc func(http.Handler) http.Handler

type httpHandlerWrapper struct {
	h     buffalo.Handler
	c     buffalo.Context
	error error
}

func (mw httpHandlerWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mw.error = mw.h(mw.c)
}

func WrapHandler(handler httpMiddlewareFunc) buffalo.MiddlewareFunc {
	return func(next buffalo.Handler) buffalo.Handler {

		wrapper := &httpHandlerWrapper{h: next}

		return func(c buffalo.Context) error {
			wrapper.c = c
			current := handler(wrapper)
			current.ServeHTTP(c.Response(), c.Request())
			return wrapper.error
		}
	}
}

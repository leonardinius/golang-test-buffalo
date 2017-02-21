package middleware

import (
	"github.com/gobuffalo/buffalo"
	"net/http"
)

type httpMiddlewareFunc func(http.Handler) http.Handler

type httpHandlerWrapper struct {
	handler buffalo.Handler
	ctx     buffalo.Context
	error   error
}

func (mw httpHandlerWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mw.error = mw.handler(mw.ctx)
}

func WrapHandler(handler httpMiddlewareFunc) buffalo.MiddlewareFunc {
	return func(next buffalo.Handler) buffalo.Handler {

		wrapper := &httpHandlerWrapper{handler: next}

		return func(context buffalo.Context) error {
			wrapper.ctx = context
			httpHandler := handler(wrapper)
			httpHandler.ServeHTTP(context.Response(), context.Request())
			return wrapper.error
		}
	}
}

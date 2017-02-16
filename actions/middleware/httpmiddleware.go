package middleware

import (
	"github.com/gobuffalo/buffalo"
	"net/http"
)

type MwHandler struct {
	h     buffalo.Handler
	c     buffalo.Context
	error error
}

func (mw MwHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mw.error = mw.h(mw.c)
}

func HttpMiddleware(handler func(http.Handler) http.Handler) func(buffalo.Handler) buffalo.Handler {
	return func(next buffalo.Handler) buffalo.Handler {

		mw := MwHandler{h: next}

		return func(c buffalo.Context) error {
			mw.c = c
			current := handler(mw)

			current.ServeHTTP(c.Response(), c.Request())
			error := mw.error

			if error != nil {
				mw.ServeHTTP(c.Response(), c.Request())
				error = mw.error
			}

			return error
		}
	}
}

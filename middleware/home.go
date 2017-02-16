package middleware

import "github.com/gobuffalo/buffalo"

// HomeHandler is a default handler to serve up
// a home page.
func APIAuthorizer(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		c.Logger().Info("APIAuthorizer::before")
		// do some work before calling the next handler
		err := next(c)
		c.Logger().Info("APIAuthorizer::after")
		// do some work after calling the next handler
		return err
	}

}

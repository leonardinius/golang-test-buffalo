package actions

import (
	"os"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/middleware"
	"github.com/leonids/test-buffalo/models"
	mw "github.com/leonids/test-buffalo/middleware"
	"github.com/markbates/going/defaults"
	"github.com/gobuffalo/buffalo/render"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = defaults.String(os.Getenv("GO_ENV"), "development")
var app *buffalo.App

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
func App() *buffalo.App {
	if app == nil {
		app = buffalo.Automatic(buffalo.Options{
			Env:         ENV,
			SessionName: "_test-buffalo_session",
		})

		app.Use(middleware.PopTransaction(models.DB))

		// index page
		app.GET("/", HomeHandler)

		g := app.Group("/api/v1")
		g.Use(mw.APIAuthorizer)
		// simple parameter tests
		g.GET("/user/", func(c buffalo.Context) error {
			name := "Hello, " + defaults.String(c.Param("name"), "<unknown>")
			return c.Render(200, render.String(name))
		})
		g.GET("/users/{name}", func(c buffalo.Context) error {
			name := "Hello, " + c.Param("name")
			return c.Render(200, render.String(name))
		})

		app.ServeFiles("/assets", assetsPath())
	}

	return app
}

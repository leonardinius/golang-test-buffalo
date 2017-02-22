package actions

import (
	"os"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/middleware"
	"github.com/gobuffalo/buffalo/render"
	"github.com/goji/httpauth"
	mw "github.com/leonids/test-buffalo/actions/middleware"
	"github.com/leonids/test-buffalo/actions/auth"
	"github.com/leonids/test-buffalo/models"
	"github.com/markbates/going/defaults"
	"gopkg.in/authboss.v1"
	"log"
	"net/http"
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

		initRoutes(app)

		app.ServeFiles("/assets", assetsPath())
	}

	return app
}

func initRoutes(app *buffalo.App) {
	// index page
	app.GET("/", HomeHandler)

	app.Resource("/users", UsersResource{&buffalo.BaseResource{}})

	{
		g := app.Group("/api/v1")
		g.Use(mw.APIAuthorizer)
		g.Use(mw.WrapHandler(httpauth.SimpleBasicAuth("leonids", "maslovs")))

		// simple parameter tests
		g.GET("/username/", func(c buffalo.Context) error {
			name := "Hello, " + defaults.String(c.Param("name"), "<unknown>")
			return c.Render(200, render.String(name))
		})
		g.GET("/username/{name}", func(c buffalo.Context) error {
			name := "Hello, " + c.Param("name")
			return c.Render(200, render.String(name))
		})
	}

	{
		g := app.Group("/api/v2")

		database := store.NewMemStorer()

		ab := authboss.New() // Usually store this globally
		ab.MountPath = "/auth"
		ab.Storer = database
		ab.OAuth2Storer = database
		ab.RootURL = `http://localhost:3000`
		ab.LogWriter = os.Stderr

		ab.XSRFName = "csrf_token"
		ab.XSRFMaker = func(_ http.ResponseWriter, r *http.Request) string {
			return "token_use_nosurf"
		}

		ab.CookieStoreMaker = store.NewCookieStorer
		ab.SessionStoreMaker = store.NewSessionStorer

		ab.Mailer = authboss.LogMailer(os.Stdout)

		ab.Policies = []authboss.Validator{
			authboss.Rules{
				FieldName:       "email",
				Required:        true,
				AllowWhitespace: false,
			},
			authboss.Rules{
				FieldName:       "password",
				Required:        true,
				MinLength:       4,
				MaxLength:       8,
				AllowWhitespace: false,
			},
		}

		if err := ab.Init(); err != nil {
			// Handle error, don't let program continue to run
			log.Fatalln(err)
		}

		// Make sure to put authboss's router somewhere
		handler := buffalo.WrapHandler(ab.NewRouter())
		g.ANY("/auth", handler)
	}
}

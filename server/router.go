package server

import (
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) Router() http.Handler {
	app := echo.New()
	app.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} ${host} '${method}' ${uri} ${latency_human} (in: ${bytes_in}b, out: ${bytes_out}b) ${error}\n",
	}))
	s.APIV1(app)
	return app
}

func (s *Server) APIV1(app *echo.Echo) *echo.Group {
	apiV1 := app.Group("/api/v1")
	apiV1.Use(session.Middleware(s.store))

	// Auth handlers
	auth := apiV1.Group("/auth")
	auth.GET("/get-session", s.GetSessionHandler)

	auth.POST("/login", s.LoginHandler)
	auth.POST("/create", s.RegisterHandler)

	auth.GET("/health", s.AuthMiddleware(func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{
			"everything": "works",
		})
	}))

	return apiV1
}

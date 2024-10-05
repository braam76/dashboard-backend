package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) Router() http.Handler {
    app := echo.New()
    s.APIV1(app)
    return app
}

func (s *Server) APIV1(app *echo.Echo) *echo.Group {
    apiV1 := app.Group("/api/v1", middleware.LoggerWithConfig(middleware.LoggerConfig{
        Format: "${time_unix} ${host} ${uri} '${method}' ${latency_human} (in: ${bytes_in}b, out: ${bytes_out}b) ${error}\n",
    }))

	// Auth handlers
    auth := apiV1.Group("/auth")
    auth.GET("/health", func(c echo.Context) error {
        return c.JSON(http.StatusOK, echo.Map{
            "everything": "works",
        })
    })
    auth.POST("/login", s.LoginHandler)

    return apiV1
}
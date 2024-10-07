package server

import (
	"log"
	"net/http"

	"github.com/braam76/dashboard-backend/server/models"
	
	"github.com/labstack/echo/v4"
)

func (s *Server) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userModel := new(models.User)

		sess, err := s.store.Get(c.Request(), "session")
		if err != nil {
			log.Fatalf("ERROR: %s", err)
		}

		if sess.Values["userID"] == nil {
			return c.String(http.StatusUnauthorized, "Unauthorized!")
		}

		result := s.db.Model(&models.User{}).
			Where("id = ?", sess.Values["userID"]).
			First(&userModel)

		if result.Error != nil {
			log.Printf("ERROR while searching for userID=%d user: %s", sess.Values["userID"], err)
			return c.String(http.StatusInternalServerError, "ERROR while searching for userID=%d user")
		}

		return next(c)
	}
}

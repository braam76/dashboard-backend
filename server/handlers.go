package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/braam76/dashboard-backend/server/models"

	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

type UserDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u UserDTO) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Username, validation.Required, validation.Length(10, 10), is.Digit),
		validation.Field(&u.Password, validation.Required),
	)
}

func (s *Server) RegisterHandler(c echo.Context) error {
	userDTO := new(UserDTO)
	userModel := new(models.User)

	if err := c.Bind(userDTO); err != nil {
		log.Printf("ERROR while binding request to UserDTO: %s", err)
		return c.String(http.StatusBadRequest, "Request is missing something, or it's bad itself")
	}

	if err := userDTO.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err,
		})
	}

	result := s.db.Model(&models.User{}).
		Where(&models.User{
			Username: userDTO.Username,
			Password: userDTO.Password,
		}).
		First(&userModel)

	if result.Error == nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("User '%s' already exists", userModel.Username))
	}

	result = s.db.Model(&models.User{}).
		Create(&models.User{
			Username: userDTO.Username,
			Password: userDTO.Password,
		})

	if result.Error != nil {
		log.Printf("ERROR: %s\n", result.Error)
		return c.String(http.StatusInternalServerError, "ERROR")
	}

	return c.String(http.StatusCreated, "GOOD! CREATE")
}

func (s *Server) LoginHandler(c echo.Context) error {
	userDTO := new(UserDTO)
	userModel := new(models.User)

	if err := c.Bind(userDTO); err != nil {
		log.Printf("ERROR while binding request to UserDTO: %s", err)
		return c.String(http.StatusBadRequest, "Request is missing something, or it's bad itself")
	}

	if err := userDTO.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err,
		})
	}

	result := s.db.Model(&models.User{}).
		Where(&models.User{
			Username: userDTO.Username,
			Password: userDTO.Password,
		}).
		First(&userModel)

	if result.Error != nil {
		log.Printf("ERROR: %s\n", result.Error)
		return c.String(http.StatusInternalServerError, "ERROR")
	}

	sess, err := s.store.Get(c.Request(), "session")
	if err != nil {
		return c.String(http.StatusInternalServerError, "Cannot GET session")
	}

	sess.Values["userID"] = userModel.ID
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return err
	}

	log.Printf("Session ID: %s\n", sess.ID)
	log.Printf("User ID from Session: %v\n", sess.Values["userID"])

	return c.JSON(http.StatusOK, sess.Values["userID"])
}

func (s *Server) GetSessionHandler(c echo.Context) error {
	sess, err := s.store.Get(c.Request(), "session")
	if err != nil {
		log.Printf("ERROR: %s\n", err)
		return c.String(http.StatusInternalServerError, "ERROR while creating session")
	}

	sess.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
	}

	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return err
	}

	log.Printf("Session ID: %s\n", sess.ID)
	log.Printf("User ID from Session: %v\n", sess.Values["userID"])

	return c.NoContent(http.StatusOK)
}

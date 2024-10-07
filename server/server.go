package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/braam76/dashboard-backend/config"
	"github.com/braam76/dashboard-backend/server/models"

	"github.com/gorilla/sessions"
	"gorm.io/gorm"

	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port  int
	db    *gorm.DB
	store *sessions.CookieStore
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("SERVER_PORT"))
	NewServer := &Server{
		port:  port,
		db:    config.NewDB(),
		store: config.NewStore(),
	}

	// Database models
	NewServer.db.AutoMigrate(
		models.User{},
	)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: NewServer.Router(),
	}

	return server
}

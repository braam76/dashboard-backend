package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"gorm.io/gorm"

	"github.com/braam76/dashboard-backend/config"

	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port int
	db   *gorm.DB
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("SERVER_PORT"))
	NewServer := &Server{
		port: port,
		db:   config.NewDB(),
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: NewServer.Router(),
	}

	return server
}

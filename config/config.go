package config

import (
	"fmt"
	"log"
	"os"

	"github.com/gorilla/sessions"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	_ "github.com/joho/godotenv/autoload"
)

var (
	addr       = os.Getenv("MYSQL_ADDR")
	username   = os.Getenv("MYSQL_USER")
	password   = os.Getenv("MYSQL_PASSWORD")
	database   = os.Getenv("MYSQL_DATABASE")
	sessionKey = os.Getenv("SESSION_KEY")
)

func NewDB() *gorm.DB {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, addr, database,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func NewStore() *sessions.CookieStore {
	return sessions.NewCookieStore([]byte(sessionKey))
}

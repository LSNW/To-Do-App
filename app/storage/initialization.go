package storage

import (
  	"gorm.io/driver/postgres"
	"github.com/gofiber/storage/redis"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"
	"os"
)

var DB *gorm.DB
var Store *session.Store

func InitDatabase() {
	dsn := "host=localhost user=todo password=" + os.Getenv("TODOPASSWORD") + " dbname=todo port=5432 sslmode=disable"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
}

func InitCache() {
	rdb := redis.New(redis.Config{
		Port: 6379,
	})
	Store = session.New(session.Config{
		Storage: rdb,
	})
}
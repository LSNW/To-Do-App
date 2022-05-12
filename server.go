package main

import (
	"github.com/gofiber/fiber/v2"
  	"gorm.io/driver/postgres"
 	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	login string
	password string
}

type ToDo struct {
	gorm.Model
	task string
	assignment string
	status string
	delete string
}

func main() {
	dsn := "host=localhost user=todo password=todopassword dbname=todo port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("")
	}
	db.AutoMigrate(&User{})
	db.AutoMigrate(&ToDo{})
	
	db.Create(&User{login: "john", password: "john_password"})
	db.Create(&ToDo{task: "exTask", assignment: "exAssignment", status:"exStatus", delete:"exDelete"})

	var todo ToDo
	db.First(&todo, 1)

  	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
    	return fiber.NewError(782, "Custom error message")
	})

  	app.Listen(":3000")
}
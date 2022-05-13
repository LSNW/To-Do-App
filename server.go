package main

import (
	"github.com/gofiber/fiber/v2"
  	"gorm.io/driver/postgres"
 	"gorm.io/gorm"
	//"fmt"
)

type User struct {
	gorm.Model
	Login string
	Password string
}

type ToDo struct {
	gorm.Model
	Task string
	Assignment string
	Status string
	Delete string
}

func main() {
	dsn := "host=localhost user=todo password=todopassword dbname=todo port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("")
	}

	db.AutoMigrate(&User{})
	db.AutoMigrate(&ToDo{})
	
	db.Create(&User{Login: "john", Password: "john_password"})
	db.Create(&ToDo{Task: "exTask", Assignment: "exAssignment", Status:"exStatus", Delete:"exDelete"})


  	app := fiber.New()

	app.Get("/:task", func(c *fiber.Ctx) error {
		var todo ToDo
		result :=  db.First(&todo, "task = ?", c.Params("task"))

		if result.RowsAffected == 0 {
			return c.SendStatus(404)
		}
	
		return c.Status(200).JSON(todo)
  	})

  	app.Listen(":3000")
}
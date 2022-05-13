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
	Task string `json:"task"`
	Assignment string `json:"assignment"`
	Status string `json:"status"`
	Delete string `json:"delete"`
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

  	app := fiber.New()

	app.Post("/newtask/", func(c *fiber.Ctx) error {
		var todo ToDo

		if err := c.BodyParser(&todo); err != nil {
			return err
		}

		db.Create(&ToDo{Task: todo.Task, Assignment: todo.Assignment, Status: todo.Status, Delete: todo.Delete})
		return c.SendString("Successfully created")
	})
	
	app.Get("/search/:task", func(c *fiber.Ctx) error {
		var todo ToDo
		result :=  db.Last(&todo, "task = ?", c.Params("task"))

		if result.RowsAffected == 0 {
			return c.SendStatus(404)
		}
	
		return c.Status(200).JSON(todo)
  	})

	app.Put("/update/:task", func(c *fiber.Ctx) error {
		var todo ToDo
		var updatedToDo ToDo

		db.Last(&todo, "task = ?", c.Params("task"))
		if err := c.BodyParser(&updatedToDo); err != nil {
			return err
		}

		if updatedToDo.Task != "" {
			todo.Task = updatedToDo.Task
		}
		if updatedToDo.Assignment != "" {
			todo.Assignment = updatedToDo.Assignment
		}
		if updatedToDo.Status != "" {
			todo.Status = updatedToDo.Status
		}
		if updatedToDo.Delete != "" {
			todo.Delete = updatedToDo.Delete
		}

		db.Save(&todo)
		return c.SendString("Successfully updated")
	})


  	app.Listen(":3000")
}
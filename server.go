package main

import (
	"github.com/gofiber/fiber/v2"
  	"gorm.io/driver/postgres"
	"github.com/gofiber/storage/redis"
	"gorm.io/gorm"
	"strconv"
	"github.com/gofiber/fiber/v2/middleware/session"
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
	
  	app := fiber.New()

	rdb := redis.New(redis.Config{
		Port: 6379,
	})

	store := session.New(session.Config{
		Storage: rdb,
	})

	_ = store

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

	app.Delete("/delete/:task", func(c *fiber.Ctx) error {
		var todo ToDo
		result := db.Where("task = ?", c.Params("task")).Delete(&todo)
		//return c.Status(200).JSON(todo)
		return c.SendString("Deleted " + strconv.Itoa(int(result.RowsAffected)) + " entries with " + c.Params("task"))
  	})


  	app.Listen(":3000")
}
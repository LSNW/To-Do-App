package main

import (
	"github.com/gofiber/fiber/v2"
	//"fmt

	"ToDoApp/app/models"
	"ToDoApp/app/controllers"
	"ToDoApp/app/storage"
)


func main() {
	app := fiber.New()
	storage.InitDatabase()
	storage.InitCache()
	storage.DB.AutoMigrate(&models.User{}, &models.ToDo{})
	storage.DB.AutoMigrate(&models.User{}, &models.ToDo{})

	// Logins
	app.Post("/login", controllers.Login)
	app.Get("/login", controllers.AutoLogin)

	// this landing page is for testing session authenticity
	app.Get("/", controllers.Landing)

	// User REST API
	app.Post("/api/User/", controllers.CreateUser)
	app.Get("/api/user/:user", controllers.FindUser)


	// ToDo REST API
	app.Post("/api/ToDo/", controllers.CreateToDo)
	app.Get("/api/ToDo/all", controllers.GetToDo)
	app.Get("/api/ToDo/:id", controllers.FindToDo)
	app.Patch("/api/ToDo/:id", controllers.UpdateToDo)
	app.Delete("/api/ToDo/:id", controllers.DeleteToDo)


  	app.Listen(":3000")
}
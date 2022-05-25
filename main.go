package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	//"fmt

	"ToDoApp/app/models"
	"ToDoApp/app/controllers"
	"ToDoApp/app/storage"
)


func main() {
	engine := html.New("./resources/views", ".html")
    app := fiber.New(fiber.Config{
        Views: engine,
    })
	storage.InitDatabase()
	storage.InitCache()
	storage.DB.AutoMigrate(&models.User{}, &models.ToDo{})

	// Logins
	app.Get("/login", controllers.Login)
	app.Post("/authenticate", controllers.Authenticate)
	//app.Get("/login", controllers.AutoLogin)

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
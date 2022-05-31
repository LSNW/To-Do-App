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
	engine := html.New("./public/views", ".html")
    app := fiber.New(fiber.Config{
        Views: engine,
    })
	app.Static("/", "./public")
	storage.InitDatabase()
	storage.InitCache()
	storage.DB.AutoMigrate(&models.User{}, &models.ToDo{})

	// Logins
	app.Get("/login", controllers.Login)
	app.Post("/login/auth", controllers.Authenticate)
	//app.Get("/login", controllers.AutoLogin)

	// this is the homepage
	app.Get("/", controllers.Landing)

	// User REST API
	app.Post("/api/User/", controllers.CreateUser)
	app.Get("/api/user/:user", controllers.FindUser)


	// ToDo REST API
	app.Post("/api/ToDo/", controllers.CreateToDo)
	app.Get("/api/ToDo/:id", controllers.FindToDo)
	app.Patch("/api/ToDo/:id", controllers.UpdateToDo)
	app.Delete("/api/ToDo/:id", controllers.DeleteToDo)

  	app.Listen(":3000")
}
package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"

	"ToDoApp/app/models"
	"ToDoApp/app/controllers"
	"ToDoApp/app/storage"
)

func initApp() *fiber.App {
	engine := html.New("./public/views", ".html")
    app := fiber.New(fiber.Config{
        Views: engine,
    })
	app.Static("/", "./public")
	storage.InitDatabase()
	storage.InitCache()
	storage.DB.AutoMigrate(&models.User{}, &models.ToDo{})

	// Logins
	app.Get("/signup", controllers.SignUp)
	app.Get("/login", controllers.Login)
	app.Get("/signout", controllers.Signout)
	app.Post("/login/auth", controllers.Authenticate)

	// Homepage
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

	return app
}

func main() {
	app := initApp()
	_ = app
}
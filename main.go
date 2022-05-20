package main

import (
	"github.com/gofiber/fiber/v2"
  	"gorm.io/driver/postgres"
	"github.com/gofiber/storage/redis"
	"gorm.io/gorm"
	"strconv"
	"github.com/gofiber/fiber/v2/middleware/session"
	"time"
	"github.com/alexedwards/argon2id"
	"github.com/jinzhu/copier"
	//"fmt

	"ToDoApp/app/models"
)


func main() {
	// app, database, cache initializations
	dsn := "host=localhost user=todo password=todopassword dbname=todo port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("")
	}
	db.AutoMigrate(&models.User{}, &models.ToDo{})
  	app := fiber.New()
	rdb := redis.New(redis.Config{
		Port: 6379,
	})
	store := session.New(session.Config{
		Storage: rdb,
	})

	expirationTime := 15 * time.Minute

	resetCookieExpiration := func(c *fiber.Ctx) {
		cookie := new(fiber.Cookie)
		cookie.Name = "session_id"
		cookie.Value = c.Cookies("session_id")
		cookie.Expires = time.Now().Add(expirationTime)
		c.Cookie(cookie)
	}

	// this is a proper login
	app.Post("/login", func(c *fiber.Ctx) error {
		var attemptedUser models.User
		var existsUser models.User

		if err := c.BodyParser(&attemptedUser); err != nil {
			return err
		}
		result := db.Last(&existsUser, "login = ?", attemptedUser.Login)

		if attemptedUser.Login == "" || attemptedUser.Password == "" {
			return c.SendString("Please enter a login and password")
		}
		match, err := argon2id.ComparePasswordAndHash(attemptedUser.Password, existsUser.Password)
		if err != nil {
			return err
		} else if !match || result.RowsAffected == 0 {
			// this msg is for debugging
			return c.SendString("Incorrect login/password")
		}

		var responseUser models.UserResponse
		copier.Copy(&responseUser, &existsUser)
		sess, err := store.Get(c)
		if err != nil {
			panic(err)
		}
		sess.Regenerate()
		defer func() {
			sess.Set("user_id", responseUser.ID)
			if err := sess.Save(); err != nil {
				panic(err)
			}
		}()

		return c.SendString(existsUser.Login +  ", you are now logged in")
	})

	// this is an automatic login
	app.Get("/login", func(c *fiber.Ctx) error {
		sess, err := store.Get(c)
		if err != nil {
			panic(err)
		}
		sess.Regenerate()
		defer func() {
			sess.Set("user_id", 154) // 154 is just a random user id
			if err := sess.Save(); err != nil {
				panic(err)
			}
		}()

		return c.SendString("You are now \"logged in\"")
	})

	// this landing page is for testing session authenticity
	app.Get("/", func(c *fiber.Ctx) error {
		sess, err := store.Get(c)
		if err != nil {
			panic(err)
		}
		if c.Cookies("session_id") != sess.ID() {
			return c.SendStatus(401)
		}
		resetCookieExpiration(c)
		return c.SendString("Welcome, your user ID is " +  strconv.Itoa(int(sess.Get("user_id").(uint))))
	})

	// User REST API
	app.Post("/api/User/", func(c *fiber.Ctx) error {
		var user models.User
		if err := c.BodyParser(&user); err != nil {
			return err
		}
		// separation is for debug
		if user.Login == "" || user.Password == "" {
			return c.SendString("Please enter a login and password")
		} else if db.Last(&user, "login = ?", user.Login).RowsAffected > 0 {
			return c.SendString("Login already exists")
		}
		hash, err := argon2id.CreateHash(user.Password, argon2id.DefaultParams)
		user.Password = hash
		if err != nil {
			return err
		}
		db.Create(&models.User{Login: user.Login, Password:user.Password})
		return c.SendString("Successfully created user profile for  " + user.Login)
	})

	app.Get("/api/user/:user", func(c *fiber.Ctx) error {
		var user models.User
		var userResponse models.UserResponse
		result := db.Last(&user, "login = ?", c.Params("user"))
		if result.RowsAffected == 0 {
			return c.SendStatus(404)
		}
		copier.Copy(&userResponse, &user)

		return c.Status(200).JSON(userResponse)
	})


	// ToDo REST API
	app.Post("/api/ToDo/", func(c *fiber.Ctx) error {
		var todo models.ToDo
		sess, err := store.Get(c)
		if err != nil {
			panic(err)
		} else if c.Cookies("session_id") != sess.ID() {
			return c.SendStatus(401)
		}

		if err := c.BodyParser(&todo); err != nil {
			return err
		}
		createToDo := models.ToDo{Task: todo.Task, Assignment: todo.Assignment, 
			Status: todo.Status, Delete: todo.Delete, UserID: sess.Get("user_id").(uint)}
		db.Create(&createToDo)
		return c.Status(200).JSON(createToDo)
	})
	
	app.Get("/api/ToDo/all", func(c *fiber.Ctx) error {
		var todos []models.ToDo
		sess, err := store.Get(c)
		if err != nil {
			panic(err)
		} else if c.Cookies("session_id") != sess.ID() {
			return c.SendStatus(401)
		}
		db.Where(&models.ToDo{UserID: sess.Get("user_id").(uint)}).Find(&todos)
		
		return c.Status(200).JSON(todos)
  	})

	app.Get("/api/ToDo/:id", func(c *fiber.Ctx) error {
		var todo models.ToDo
		sess, err := store.Get(c)

		result :=  db.Last(&todo, c.Params("id"))

		if err != nil {
			panic(err)
		} else if c.Cookies("session_id") != sess.ID() || todo.UserID != sess.Get("user_id").(uint) {
			return c.SendStatus(401)
		} else if result.RowsAffected == 0 {
			return c.SendStatus(404)
		}
		
		return c.Status(200).JSON(todo)
  	})

	app.Patch("/api/ToDo/:id", func(c *fiber.Ctx) error {
		var todo models.ToDo
		var updatedToDo models.ToDo
		sess, err := store.Get(c)
		if err != nil {
			return err
		}

		result := db.Last(&todo, c.Params("id"))

		if err := c.BodyParser(&updatedToDo); err != nil {
			return err
		} else if c.Cookies("session_id") != sess.ID() || todo.UserID != sess.Get("user_id").(uint) {
			return c.SendStatus(401)
		} else if result.RowsAffected == 0 {
			return c.SendStatus(404)
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

	app.Delete("/api/ToDo/:id", func(c *fiber.Ctx) error {
		var todo models.ToDo
		sess, err := store.Get(c)
		if err != nil {
			return err
		} else if c.Cookies("session_id") != sess.ID() || todo.UserID != sess.Get("user_id").(uint) {
			return c.SendStatus(401)
		}
		result := db.First(&todo, c.Params("task")).Delete(&todo)
		if result.RowsAffected == 0 {
			return c.SendStatus(404)
		}
		return c.SendString("Task has been deleted")
  	})


  	app.Listen(":3000")
}
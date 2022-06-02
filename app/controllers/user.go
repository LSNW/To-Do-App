package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/alexedwards/argon2id"
	"github.com/jinzhu/copier"
	"strconv"
	"log"

	"ToDoApp/app/models"
	"ToDoApp/app/storage"
)

func getToDo(c *fiber.Ctx) []models.ToDo {
	var todos []models.ToDo
	sess, err := storage.Store.Get(c)
	if err != nil {
		log.Println(err)
	}
	storage.DB.Where(&models.ToDo{UserID: sess.Get("user_id").(uint)}).Order("id asc").Find(&todos)
	
	return todos
}

func Login(c *fiber.Ctx) error {
	return c.Render("login", nil)
}

func Authenticate(c *fiber.Ctx) error {
	attemptedUser := models.User{Login:c.FormValue("login"), Password:c.FormValue("password")}
	var existsUser models.User

	result := storage.DB.Last(&existsUser, "login = ?", attemptedUser.Login)

	if attemptedUser.Login == "" || attemptedUser.Password == "" {
		return c.Render("login", fiber.Map {
			"error": "Please enter a username and password",
		})
	}

	match, err := argon2id.ComparePasswordAndHash(attemptedUser.Password, existsUser.Password)
	if !match || result.RowsAffected == 0 {
		return c.Render("login", fiber.Map {
			"error": "Incorrect username or password",
		})
	} else if err != nil {
		log.Println(err)
		return c.SendStatus(500)
	} 

	var responseUser models.UserResponse
	copier.Copy(&responseUser, &existsUser)
	sess, err := storage.Store.Get(c)
	if err != nil {
		log.Println(err)
		return c.SendStatus(500)
	}
	sess.Regenerate()
	defer func() {
		sess.Set("user_id", responseUser.ID)
		if err := sess.Save(); err != nil {
			log.Println(err)
			//return c.SendStatus(500)
		}
	}()

	return c.Redirect("/")
}

func AutoLogin(c *fiber.Ctx) error {
	sess, err := storage.Store.Get(c)
	if err != nil {
		log.Println(err)
		return c.SendStatus(500)
	}
	sess.Regenerate()
	defer func() {
		sess.Set("user_id", uint(154)) // 154 is just a random user id
		if err := sess.Save(); err != nil {
			log.Println(err)
			//return c.SendStatus(500)
		}
	}()

	return c.SendString("You are now \"logged in\"")
}

func Landing(c *fiber.Ctx) error {
	sess, err := storage.Store.Get(c)
	if err != nil {
		log.Println(err)
		return c.SendStatus(500)
	}
	if c.Cookies("session_id") != sess.ID() {
		return c.Redirect("/login")
	}
	return c.Render("index", fiber.Map{
		"user_id": strconv.Itoa(int(sess.Get("user_id").(uint))),
		"todos": getToDo(c),
	})
}

func CreateUser(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		log.Println(err)
		return c.SendStatus(500)
	}
	// separation is for debug
	if user.Login == "" || user.Password == "" {
		return c.SendString("Please enter a login and password")
	} else if storage.DB.Last(&user, "login = ?", user.Login).RowsAffected > 0 {
		return c.SendString("Login already exists")
	}
	hash, err := argon2id.CreateHash(user.Password, argon2id.DefaultParams)
	user.Password = hash
	if err != nil {
		log.Println(err)
		return c.SendStatus(500)
	}
	storage.DB.Create(&models.User{Login: user.Login, Password:user.Password})
	return c.SendString("Successfully created user profile for  " + user.Login)
}

func FindUser(c *fiber.Ctx) error {
	var user models.User
	var userResponse models.UserResponse
	result := storage.DB.Last(&user, "login = ?", c.Params("user"))
	if result.RowsAffected == 0 {
		return c.SendStatus(404)
	}
	copier.Copy(&userResponse, &user)

	return c.Status(200).JSON(userResponse)
}
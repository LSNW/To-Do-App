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

var message string

func getToDo(c *fiber.Ctx) []models.ToDo {
	var todos []models.ToDo
	sess, err := storage.Store.Get(c)
	if err != nil {
		log.Println(err)
	}
	storage.DB.Where(&models.ToDo{UserID: sess.Get("user_id").(uint)}).Order("id asc").Find(&todos)
	
	return todos
}

func SignUp(c *fiber.Ctx) error {
	defer func () {
		message = ""
	}()
	return c.Render("signup", fiber.Map {
		"error": message,
	})
}

func Login(c *fiber.Ctx) error {
	defer func () {
		message = ""
	}()
	return c.Render("login", fiber.Map {
		"error": message,
	})
}

func Authenticate(c *fiber.Ctx) error {
	attemptedUser := models.User{Login:c.FormValue("login"), Password:c.FormValue("password")}
	var existsUser models.User

	result := storage.DB.Last(&existsUser, "login = ?", attemptedUser.Login)

	if attemptedUser.Login == "" || attemptedUser.Password == "" {
		message = "Please enter a username and password"
		return c.Redirect("/login")
	}

	match, err := argon2id.ComparePasswordAndHash(attemptedUser.Password, existsUser.Password)
	if !match || result.RowsAffected == 0 {
		message = "Incorrect username or password"
		return c.Redirect("/login")
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
		sess.Set("login", responseUser.Login)
		if err := sess.Save(); err != nil {
			log.Println(err)
			//return c.SendStatus(500)
		}
	}()

	return c.Redirect("/")
}

func Signout(c *fiber.Ctx) error {
	sess, err := storage.Store.Get(c)
	if err != nil {
		log.Println(err)
		return c.SendStatus(500)
	}
	if err := sess.Destroy(); err != nil {
		log.Println(err)
	}
	return c.Redirect("/login")
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
		"login": sess.Get("login"),
		"todos": getToDo(c),
	})
}

func CreateUser(c *fiber.Ctx) error {
	user := models.User{Login:c.FormValue("login"), Password:c.FormValue("password")}

	// separation is for debug
	if user.Login == "" || user.Password == "" {
		return c.SendString("Please enter a login and password")
	} else if storage.DB.Last(&user, "login = ?", user.Login).RowsAffected > 0 {
		message = "Login already exists"
		return c.Redirect("/signup")
	}
	hash, err := argon2id.CreateHash(user.Password, argon2id.DefaultParams)
	user.Password = hash
	if err != nil {
		log.Println(err)
		return c.SendStatus(500)
	}
	storage.DB.Create(&models.User{Login: user.Login, Password:user.Password})
	return c.Redirect("/login")
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
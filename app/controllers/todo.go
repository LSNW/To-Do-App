package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"log"
	"time"

	"ToDoApp/app/models"
	"ToDoApp/app/storage"
)

var expirationTime = 60 * time.Minute

func resetCookieExpiration(c *fiber.Ctx) {
	cookie := new(fiber.Cookie)
	cookie.Name = "session_id"
	cookie.Value = c.Cookies("session_id")
	cookie.Expires = time.Now().Add(expirationTime)
	c.Cookie(cookie)
}

func CreateToDo(c *fiber.Ctx) error {
	var todo models.ToDo
	sess, err := storage.Store.Get(c)
	if err != nil {
		log.Println(err)
		return c.SendStatus(500)
	} else if c.Cookies("session_id") != sess.ID() {
		return c.SendStatus(401)
	}

	if err := c.BodyParser(&todo); err != nil {
		log.Println(err)
		return c.SendStatus(500)
	}
	createToDo := models.ToDo{Task: todo.Task, Assignment: todo.Assignment, 
		Status: todo.Status, Delete: todo.Delete, UserID: sess.Get("user_id").(uint)}
	storage.DB.Create(&createToDo)
	return c.Status(200).JSON(createToDo)
}

func GetToDo(c *fiber.Ctx) error {
	var todos []models.ToDo
	var todoDTOs []models.ToDoDTO
	sess, err := storage.Store.Get(c)
	if err != nil {
		log.Println(err)
		return c.SendStatus(500)
	} else if c.Cookies("session_id") != sess.ID() {
		return c.SendStatus(401)
	}
	storage.DB.Where(&models.ToDo{UserID: sess.Get("user_id").(uint)}).Find(&todos)
	copier.Copy(&todoDTOs, &todos)
	return c.Status(200).JSON(todos)
}

func FindToDo(c *fiber.Ctx) error {
	var todo models.ToDo
	sess, err := storage.Store.Get(c)

	result :=  storage.DB.Last(&todo, c.Params("id"))

	if err != nil {
		log.Println(err)
		return c.SendStatus(500)
	} else if c.Cookies("session_id") != sess.ID() || todo.UserID != sess.Get("user_id").(uint) {
		return c.SendStatus(401)
	} else if result.RowsAffected == 0 {
		return c.SendStatus(404)
	}
	
	return c.Status(200).JSON(todo)
  }

func UpdateToDo(c *fiber.Ctx) error {
	var todo models.ToDo
	var updatedToDo models.ToDo
	sess, err := storage.Store.Get(c)
	if err != nil {
		log.Println(err)
		return c.SendStatus(500)
	}

	result := storage.DB.Last(&todo, c.Params("id"))

	if err := c.BodyParser(&updatedToDo); err != nil {
		log.Println(err)
		return c.SendStatus(500)
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

	storage.DB.Save(&todo)
	return c.Status(200).JSON(todo)
}

func DeleteToDo(c *fiber.Ctx) error {
	var todo models.ToDo
	sess, err := storage.Store.Get(c)
	result := storage.DB.First(&todo, c.Params("id"))
	if err != nil {
		log.Println(err)
		return c.SendStatus(500)
	} else if c.Cookies("session_id") != sess.ID() || todo.UserID != sess.Get("user_id").(uint) {
		return c.SendStatus(401)
	} else if result.RowsAffected == 0 {
		return c.SendStatus(404)
	}
	result.Delete(&todo)
	return c.SendStatus(200)
  }
package models

import (
	"gorm.io/gorm"
)

type ToDo struct {
	gorm.Model
	Task string `json:"task"`
	Assignment string `json:"assignment"`
	Status string `json:"status"`
	Delete string `json:"delete"`
	UserID uint
}

type ToDoResponse struct {
	ID uint `copier:"must, nopanic"`
	Task string `copier:"must, nopanic"`
	Assignment string `copier:"must, nopanic"`
	Status string `copier:"must, nopanic"`
	Delete string `copier:"must, nopanic"`
}
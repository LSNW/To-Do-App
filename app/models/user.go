package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Login string `json:"login"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID uint `copier:"must"`
	Login string `copier:"must"`
}
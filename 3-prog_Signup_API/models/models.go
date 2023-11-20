package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name         string
	Email        string
	PasswordHash string
}

//take request body from user 
type NewUserReq struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

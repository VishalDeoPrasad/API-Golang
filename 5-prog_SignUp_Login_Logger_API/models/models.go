package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name         string
	Email        string
	PasswordHash string
}

//take request body from user for new record
type NewUserReq struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

//take the request body from user for login
type LoginReq struct{
	Email string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

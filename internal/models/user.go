package models

import "gorm.io/gorm"

type User struct {
    gorm.Model
	FirstName string `json:"firstname" validate:"required"`
    Surname   string `json:"surname" validate:"required"`
    Username  string `json:"username" gorm:"unique;not null"`
    Email     string `json:"email" gorm:"unique;not null" validate:"required,email"`
	Password  string `json:"password" gorm:"not null" validate:"required,min=6"`
}

type UserLogin struct {
    Email    string `json:"email" validate:"required"`
    Password string `json:"password" validate:"required"`
}

package models

import (
    "github.com/google/uuid"
    "time"
)

type User struct {
    ID        string    `json:"id"`
	FirstName string    `json:"firstname" validate:"required"`
    Surname   string    `json:"surname" validate:"required"`
    Username  string    `json:"username"`
    Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"password" validate:"required,min=6"`
    IsAdmin   bool      `json:"is_admin"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) AddDefaults() {
    u.ID = uuid.New().String()
    u.IsAdmin = false
    u.CreatedAt = time.Now()
    u.UpdatedAt = time.Now()

}

type UserLogin struct {
    Email    string `json:"email" validate:"required"`
    Password string `json:"password" validate:"required"`
}

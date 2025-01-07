package models

import (
	"time"

	"github.com/google/uuid"
	// "github.com/google/uuid"
	// "github.com/blacac3/go-rest-api/internal/util"
)

var err error

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
    u.ID = GenerateUUID(u.Email)
    u.IsAdmin = false
    u.CreatedAt = time.Now()
    u.UpdatedAt = time.Now()
}



type UserLogin struct {
    Email    string `json:"email" validate:"required"`
    Password string `json:"password" validate:"required"`
}

// GenerateUUID generates a UUIDv5 using a namespace and a string (email in this case)
func GenerateUUID(email string) string {
	namespace := uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	return uuid.NewMD5(namespace, []byte(email)).String()
}




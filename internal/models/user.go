package models

// import "gorm.io/gorm"

type User struct {
    ID        string `json:"id" validate:"required"`
	FirstName string `json:"firstname" validate:"required"`
    Surname   string `json:"surname" validate:"required"`
    Username  string `json:"username"`
    Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6"`
    Role      string `json:"role"`
}

type UserLogin struct {
    Email    string `json:"email" validate:"required"`
    Password string `json:"password" validate:"required"`
}

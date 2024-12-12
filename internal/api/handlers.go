package api

import (
	// "fmt"
	"log"
	"net/http"

	"github.com/blacac3/go-rest-api/internal/models"
	"github.com/blacac3/go-rest-api/internal/util"
	"gorm.io/gorm"
)


var messagePayload map[string]interface{} = make(map[string]interface{})

func ChangeDB(db *gorm.DB){
    DB = db
}

func HandleHealthz(w http.ResponseWriter, r *http.Request){
    payload := messagePayload
    payload["message"] = "Server is Ok 👌"
    util.RespondWithJson(w, http.StatusCreated, payload)
    return
}

func HandleLogin(w http.ResponseWriter, r *http.Request){
    type LoginRequest struct{
        Email string `json:"email"`
        Password string `json:"password"`
    }
    var request LoginRequest
    var user models.User
    payload := messagePayload

    if err:=util.ValidateRequest(*r, &request); err != nil{
        payload := messagePayload
        payload["message"] = "Error while logging in"
        util.RespondWithJson(w, http.StatusBadRequest, payload)
        log.Printf("Authentication Failed: %v", err)
        return
    }
    if err := DB.First(&user, "email = ? ", request.Email).Error; err == nil {
        if verified:=util.VerifyPassword(request.Password, user.Password); verified == false{
            payload["error"] = "Invalid Password"
            util.RespondWithJson(w, http.StatusUnauthorized, payload)
            return
        }
    }else{        
        payload["error"] = "User not found"
        log.Printf("User not found: %v", err)
        util.RespondWithJson(w, http.StatusNotFound, payload)
        return
    }
    util.RespondWithJson(w, 200, user)
    return
}


func HandleRegisteration(w http.ResponseWriter, r *http.Request){
    var user models.User
    if err := util.ValidateRequest(*r, &user); err != nil{
        payload:=messagePayload
        payload["error"] = "Bad Request"
        util.RespondWithJson(w, http.StatusBadRequest, payload)
        log.Printf("User Registration Failed: %v", err)
        return
    }
    user.Password = util.HashPassword(user.Password)
    if err := DB.Create(&user).Error; err != nil{
        payload := messagePayload
        payload["error"] = "Error occured while adding user"
        if err:= DB.Where("email = ?", user.Email).First(&user).Error; err == nil{
            payload["error"] = "User already exists"
        }
        util.RespondWithJson(w, http.StatusBadRequest, payload)
        log.Printf("Error while adding user: %v", err)
        return
    }
    
    util.RespondWithJson(w, http.StatusCreated, user)
    return
}



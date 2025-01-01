package api

import (
	// "fmt"
	"log"
	"net/http"

	"github.com/blacac3/go-rest-api/internal/database"
	"github.com/blacac3/go-rest-api/internal/models"
	"github.com/blacac3/go-rest-api/internal/util"
	"github.com/gin-gonic/gin"
)

var (
    messagePayload map[string]interface{} = make(map[string]interface{})
    users_collection_name string = "users"
    
)
//Use when integrating Postgresql
// var DB =database.DB

// func ChangeDB(db *gorm.DB){
//     DB = db
// }

func HandleHealthz(c *gin.Context){
	c.JSON(http.StatusOK, gin.H{"message": "Server is up and running! 👌"})
}

func HandleLogin(c *gin.Context){
    type LoginRequest struct{
        Email string `json:"email"`
        Password string `json:"password"`
    }
    var request LoginRequest
    var user models.User
    
    if err := database.CreateBoltBucket(users_collection_name); err != nil{
        log.Print(err)
        c.JSON(http.StatusNotFound, gin.H{"error":"Server Error"})
        return
    }

    if err:=util.ValidateRequest(c, &request); err != nil{
        c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
        return
    }
    userInfo, err := database.GetBoltBucket(users_collection_name, request.Email)
    if err != nil{
        log.Print(err)
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }
    util.Deserialize(userInfo, &user)

    if verified:=util.VerifyPassword(request.Password, user.Password); verified == false{
        c.JSON(http.StatusNotFound, gin.H{"error": "Invalid Password"})
        return
    }
    token, err := util.GenerateJWT(user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while generating JWT"})
        return
    }else{
        c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})
    }
    return
}


func HandleRegisteration(c *gin.Context){
    var user models.User
    if err := util.ValidateRequest(c, &user); err != nil{
        c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
        return
    }
    user.Password = util.HashPassword(user.Password)
    serializedUser, _ := util.Serialize(user)
    database.CreateBoltBucket(users_collection_name)
    if err := database.UpdateBoltBucket(users_collection_name, user.Email, serializedUser); err != nil{
        c.JSON(http.StatusNotFound, gin.H{"error": "Unable to register user"})
        return
    }

    response := user
    response.Password = ""
    util.RespondWithJson(c, http.StatusCreated, response)
    return
}



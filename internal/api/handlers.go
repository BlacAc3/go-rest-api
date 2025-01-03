package api

import (
	// "fmt"
	"encoding/json"
	"fmt"
	"io"
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
    files_collection_name string = "files"
    index_collection_name string = "index"
    
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
    json.Unmarshal(userInfo, &user)

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
    user.AddDefaults()

    userBytes, err := json.Marshal(user)
    if err != nil{
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while encoding user data"})
        return
    }
    database.CreateBoltBucket(users_collection_name)
    if err := database.UpdateBoltBucket(users_collection_name, user.Email, userBytes); err != nil{
        c.JSON(http.StatusNotFound, gin.H{"error": "Unable to register user"})
        return
    }

    response := user
    response.Password = ""
    util.RespondWithJson(c, http.StatusCreated, response)
    return
}


func HandleFileUpload(c *gin.Context){
    var maxFileSize int = 10 * 1024 *1024
    
    file, header, err := c.Request.FormFile("file")
    if err != nil{
        c.JSON(http.StatusBadRequest, gin.H{"error": "Error while uploading file"})
        return
    }

    defer file.Close()
    
    fileData, err := io.ReadAll(file)
    if err != nil{
        c.JSON(http.StatusBadRequest, gin.H{"error": "Error while reading file"})
        return
    }
    if len(fileData) > maxFileSize{
        c.JSON(http.StatusBadRequest, gin.H{"error": "File size is too large"})
        return
    }

    fileModel := models.File{}
    fileModel.Type = header.Header.Get("Content-Type")
    fileModel.EncryptedData, _ = util.EncryptFile(fileData)
    fileModel.AddDefaults()
    fileModel.UserID, _ = util.VerifyJWT(util.GetJWT(c))
    
    // Storing File
    fileModelBytes, err := json.Marshal(fileModel)
    if err != nil{
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while encoding file data"})
        return
    }
    database.CreateBoltBucket(files_collection_name)
    database.UpdateBoltBucket(files_collection_name, fileModel.ID, fileModelBytes)
    UpdateIndex(fileModel)

    c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully!", "user": fileModel.UserID})
    return
}

func UpdateIndex(fileModel models.File)error{
    database.CreateBoltBucket(index_collection_name)
    var userFileList []string
    userFileListBytes, err := database.GetBoltBucket(index_collection_name, fileModel.UserID)
    if err != nil{
        var fileList []string
        fileList = append(fileList, fileModel.ID)
        fileListBytes, err := json.Marshal(fileList)
        if err != nil{
            return fmt.Errorf("Error updating Index: %v", err)
        }
        database.UpdateBoltBucket(index_collection_name, fileModel.UserID, fileListBytes)
        return nil
    }
    
    json.Unmarshal(userFileListBytes, &userFileList)
    userFileList = append(userFileList, fileModel.ID)
    userFileListBytes, err = json.Marshal(userFileList)
    if err != nil{
        return fmt.Errorf("Error updating Index: %v", err)
    }

    database.UpdateBoltBucket(index_collection_name, fileModel.UserID, userFileListBytes)
    return nil
}




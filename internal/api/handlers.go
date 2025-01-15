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
	messagePayload     map[string]interface{} = make(map[string]interface{})
	users_collection   string                 = "users"
	files_collection   string                 = "files"
	user_file_index    string                 = "user_file_index"
	folders_collection string                 = "folders"
	user_folder_index  string                 = "user_folder_index"
)

//Use when integrating Postgresql
// var DB =database.DB

// func ChangeDB(db *gorm.DB){
//     DB = db
// }

func HandleHealthz(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Server is up and running! 👌"})
	return
}

func HandleLogin(c *gin.Context) {
	type LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var request LoginRequest
	var user models.User

	if err := database.CreateBoltBucket(users_collection); err != nil {
		log.Print(err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Server Error"})
		return
	}

	if err := util.ValidateRequest(c, &request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}
	userInfo, err := database.GetBoltBucket(users_collection, models.GenerateUUID(request.Email))
	if err != nil {
		// log.Print(err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	json.Unmarshal(userInfo, &user)

	if verified := util.VerifyPassword(request.Password, user.Password); verified == false {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid Password"})
		return
	}
	token, err := util.GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while generating JWT"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": string(token)})
	}
	return
}

func HandleRegisteration(c *gin.Context) {
	var user models.User
	if err := util.ValidateRequest(c, &user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}
	user.Password = util.HashPassword(user.Password)
	user.AddDefaults()

	userBytes, err := json.Marshal(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while encoding user data"})
		return
	}
	database.CreateBoltBucket(users_collection)
	if err := database.UpdateBoltBucket(users_collection, user.ID, userBytes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to register user"})
		return
	}

	response := user
	response.Password = ""
	util.RespondWithJson(c, http.StatusCreated, response)
	return
}

func HandleFileUpload(c *gin.Context) {
	var maxFileSize int = 10 * 1024 * 1024

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while uploading file"})
		return
	}

	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while reading file"})
		return
	}
	if len(fileData) > maxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size is too large"})
		return
	}

	fileModel := models.File{}
	fileModel.Filename = header.Filename
	fileModel.Type = header.Header.Get("Content-Type")
	fileModel.EncryptedData, _ = util.EncryptFile(fileData)
	fileModel.AddDefaults()
	fileModel.UserID, _ = util.VerifyJWT(util.GetJWT(c))

	// Storing File
	fileModelBytes, err := json.Marshal(fileModel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while encoding file data"})
		return
	}

	err = UpdateIndex(user_file_index, fileModel.UserID, fileModel) //Indexing files to UserId
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "updating index"})
		return
	}

	//Store files
	database.CreateBoltBucket(files_collection)
	database.UpdateBoltBucket(files_collection, fileModel.ID, fileModelBytes)

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully!", "fileID": fileModel.ID, "filename": fileModel.Filename})
	return
}

func UpdateIndex[T models.File | models.Folder](indexBucketName string, userID string, data T) error {
	var id string
	switch v := any(data).(type) {
	case models.File:
		id = v.ID
	case models.Folder:
		id = v.ID
	}

	database.CreateBoltBucket(indexBucketName)
	userItemMap := make(map[string]bool)
	ItemMapBytes, err := database.GetBoltBucket(indexBucketName, userID)
	if err != nil {
		userItemMap[id] = true
		itemMapBytes, err := json.Marshal(userItemMap)
		if err != nil {
			return fmt.Errorf("Error updating Index: %v", err)
		}
		database.UpdateBoltBucket(user_file_index, userID, itemMapBytes)
		return nil
	}

	json.Unmarshal(ItemMapBytes, &userItemMap)
	userItemMap[id] = true
	ItemMapBytes, err = json.Marshal(userItemMap)
	if err != nil {
		return fmt.Errorf("Error updating Index: %v", err)
	}

	database.UpdateBoltBucket(user_file_index, userID, ItemMapBytes)
	return nil
}

func HandleFileDownload(c *gin.Context) {
	fileID := c.Param("fileID")
	userID, err := util.VerifyJWT(util.GetJWT(c))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User Not Authenticated"})
		return
	}

	//Get files indexed to a user ID
	fileMapBytes, err := database.GetBoltBucket(user_file_index, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No files found"})
		return
	}

	fileIDMap := make(map[string]bool)
	json.Unmarshal(fileMapBytes, &fileIDMap)

	_, ok := fileIDMap[fileID]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Requested file not found"})
		return
	}
	var downloadFileModel models.File
	fileModelBytes, err := database.GetBoltBucket(files_collection, fileID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	json.Unmarshal(fileModelBytes, &downloadFileModel)

	downloadFileData, err := util.DecryptFile(downloadFileModel.EncryptedData)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while downloading file"})
		return
	}
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", downloadFileModel.Filename))
	c.Header("Content-Type", downloadFileModel.Type)
	c.Header("Content-Length", fmt.Sprintf("%d", len(downloadFileData)))

	// Send the decrypted file data to the frontend
	c.Data(http.StatusOK, downloadFileModel.Type, downloadFileData)
	return
}

func HandleFilesRetrieval(c *gin.Context) {
	userID, err := util.VerifyJWT(util.GetJWT(c))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User Not Authenticated"})
		return
	}
	//Get files list indexed to a user ID
	fileMapBytes, err := database.GetBoltBucket(user_file_index, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No files found"})
		return
	}
	fileMap := make(map[string]bool)
	json.Unmarshal(fileMapBytes, &fileMap)
	var files []interface{}
	for fileID, _ := range fileMap {
		fileModelBytes, err := database.GetBoltBucket(files_collection, fileID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "No files found"})
		}
		var file models.File
		json.Unmarshal(fileModelBytes, &file)
		fileResp := map[string]string{}
		fileResp["id"] = file.ID
		fileResp["name"] = file.Filename
		files = append(files, fileResp)
	}
	response := gin.H{
		"files":      files,
		"totalCount": len(files),
	}
	c.JSON(http.StatusOK, response)
	return
}

func HandleCreateFolder(c *gin.Context) {
	userID, err := util.VerifyJWT(util.GetJWT(c))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User Not Authenticated"})
		return
	}
	var requestBody map[string]interface{}
	err = c.ShouldBindJSON(&requestBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	folder := &models.Folder{}
	folderName, ok := requestBody["folderName"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Folder Name is required"})
		return
	}
	folder.Name = folderName

	parentFolderID, ok := requestBody["parentFolderID"].(string)
	if !ok {
		folder.ParentFolderID = ""
	}
	folder.ParentFolderID = parentFolderID

	folder.AddDefaults()
	folder_bytes, err := json.Marshal(folder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while encoding"})
		return
	}

	//Link folder to user via index
	err = UpdateIndex(user_folder_index, userID, *folder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while indexing"})
	}

	// Add folder to database
	database.CreateBoltBucket(folders_collection)
	database.UpdateBoltBucket(folders_collection, folder.ID, folder_bytes)

	c.JSON(http.StatusOK, gin.H{"folder_id": folder.ID, "message": "Folder Created Successfully"})
	return
}

func HandleMoveFolder(c *gin.Context) {
	userID, _ := util.VerifyJWT(util.GetJWT(c))
	folderMap := make(map[string]bool)
	var requestBody map[string]interface{}
	err := c.ShouldBindJSON(&requestBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	folderID, curr_ok := requestBody["current_folder_id"].(string)
	targetFolderID, targ_ok := requestBody["target_folder_id"].(string)
	if !curr_ok || !targ_ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Folder ID and Target Folder ID is required"})
		return
	}

	//check if files requested are owned by the user
	folderMapBytes, err := database.GetBoltBucket(user_folder_index, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Can not move folders that do not exist"})
		return
	}
	json.Unmarshal(folderMapBytes, &folderMap)

	_, fID_ok := folderMap[folderID]
	_, tID_ok := folderMap[targetFolderID]
	if !fID_ok || !tID_ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Can not move folders that do not exist"})
		return
	}

	folder := &models.Folder{}
	folderBytes, err := database.GetBoltBucket(folders_collection, folderID)
	json.Unmarshal(folderBytes, folder)

	folder.ParentFolderID = targetFolderID
	folderBytes, err = json.Marshal(folder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while encoding"})
		return
	}
	database.UpdateBoltBucket(folders_collection, folderID, folderBytes)
	c.JSON(http.StatusOK, gin.H{"message": "Folder Moved Successfully"})
	return
}

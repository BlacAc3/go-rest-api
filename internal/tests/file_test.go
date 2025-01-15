package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var fileUploadEndpoint string = "files/upload"


func UploadFile(router *gin.Engine, token string, fileName string, fileContent []byte) (*httptest.ResponseRecorder, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", fileName)
    if err != nil{
        return nil, fmt.Errorf("Failed to create Form File: %v", err)
    }

	_, err = part.Write(fileContent)
    if err != nil{
        return nil, fmt.Errorf("Failed to write file content to form data: %v", err)
    }
	writer.Close()

	req, err := http.NewRequest(http.MethodPost, baseUrl+"/files/upload", body)
    if err != nil {
        return nil, fmt.Errorf("Failed to Create Request for health endpoint Test ---> %v", err)
    }

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+ fmt.Sprint(token))

    res := PerformRequest(router, req)
    return res, nil
    
}

//File Upload and retrieval endpoint test
func TestFileUploadAndRetrieval(t *testing.T) {
    SetupServer()
	router := GetRouter()

    registerUser(t, user3)
    resp := LoginUser(t, user3)
    loginResponse:=resp.Body.String()
    
    var loginResult map[string]string    
    err := json.Unmarshal([]byte(loginResponse), &loginResult)
    if err != nil{
        t.Errorf("Parse Error: %v", err)
    }

    token := loginResult["token"]

    //File upload
    fileName := "testfile1.txt"
	fileContent := []byte("This is a test file content")
    res, err :=UploadFile(router, token, fileName, fileContent)
    assert.NoError(t, err)
    assert.Equal(t, http.StatusOK, res.Code)

    fileID := ""
    var response map[string]string

    t.Run("TestFileDownload", func(t *testing.T){
        //TODO Split this test, then add another test for users getting another user's files via the get file endpoint
        // Step 2: Check upload response
        responseString := res.Body.String()
        err = json.Unmarshal([]byte(responseString), &response)
        fileID = response["fileID"]
        assert.NoError(t, err)
        assert.Equal(t, "File uploaded successfully!", response["message"])
        assert.Equal(t, fileName, response["filename"])

        // Step 3: Simulate file retrieval
        req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%v/files/%v", baseUrl, fileID), nil)
        req.Header.Set("Authorization", "Bearer "+ fmt.Sprint(token))
        res = PerformRequest(router, req)

        assert.Equal(t, http.StatusOK, res.Code)
        assert.Equal(t, fileContent, res.Body.Bytes())
    })

    t.Run("TestUnauthorizedFileRetrieval", func(t *testing.T){
        registerUser(t, user2)
        resp := LoginUser(t, user2)
        loginResponse:=resp.Body.String()
        
        var loginResult map[string]string    
        err := json.Unmarshal([]byte(loginResponse), &loginResult)
        if err != nil{
            t.Errorf("Parse Error: %v", err)
        }

        token := loginResult["token"]
        req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%v/files/%v", baseUrl, fileID), nil)
        req.Header.Set("Authorization", "Bearer "+ fmt.Sprint(token))
        res = PerformRequest(router, req)
        assert.Equal(t, http.StatusNotFound, res.Code)
    })

}


// Get all files Endpoint Tests
func TestGetFiles(t *testing.T) {
    SetupServer()
	router := GetRouter()

    registerUser(t, user3)
    resp := LoginUser(t, user3)
    loginResponse:=resp.Body.String()
    
    var loginResult map[string]string    
    err := json.Unmarshal([]byte(loginResponse), &loginResult)
    if err != nil{
        t.Errorf("Parse Error: %v", err)
    }

    token := loginResult["token"]

    _, err = UploadFile(router, token, "file1.txt", []byte("content of file1"))
    assert.NoError(t, err)
    UploadFile(router, token, "file2.txt", []byte("content of file2"))
    UploadFile(router, token, "file3.txt", []byte("content of file3"))


    req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%v/files", baseUrl), nil)
	req.Header.Set("Authorization", "Bearer "+ fmt.Sprint(token))
    res := PerformRequest(router, req)
    response := make(map[string]interface{})
    err = json.Unmarshal([]byte(res.Body.String()), &response)
    if err != nil{
        t.Errorf("Parse Error: %v", err)
    }

    assert.Equal(t, float64(3), response["totalCount"])
    assert.Equal(t, http.StatusOK, res.Code)
}

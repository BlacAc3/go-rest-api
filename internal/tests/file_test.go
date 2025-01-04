package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blacac3/go-rest-api/internal/api"
	"github.com/stretchr/testify/assert"
)

var fileUploadEndpoint string = "files/upload"

func TestFileUploadAndRetrieval(t *testing.T) {
    SetupServer()
	router := api.InitRouter()

    registerUser(t, user3)
    resp := LoginUser(t, user3)
    loginResponse:=resp.Body.String()
    
    var loginResult map[string]string    
    err := json.Unmarshal([]byte(loginResponse), &loginResult)
    if err != nil{
        t.Errorf("Parse Error: %v", err)
    }

    token := loginResult["token"]

	// Step 1: Simulate file upload
	fileContent := []byte("This is a test file content")
	fileName := "testfile.txt"

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", fileName)
	assert.NoError(t, err)

	_, err = part.Write(fileContent)
	assert.NoError(t, err)
	writer.Close()

	req, err := http.NewRequest(http.MethodPost, baseUrl+"/files/upload", body)
    if err != nil {
        t.Fatalf("Failed to Create Request for health endpoint Test ---> %v", err)
    }

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+ fmt.Sprint(token))


    res := PerformRequest(router, req)
	assert.Equal(t, http.StatusOK, res.Code)

	// Step 2: Check upload response
    responseString := res.Body.String()
	var response map[string]string
	err = json.Unmarshal([]byte(responseString), &response)
    fileID := response["fileID"]
	assert.NoError(t, err)
	assert.Equal(t, "File uploaded successfully!", response["message"])
	assert.Equal(t, fileName, response["filename"])

	// Step 3: Simulate file retrieval
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("%v/files/%v", baseUrl, fileID), nil)
	req.Header.Set("Authorization", "Bearer "+ fmt.Sprint(token))
    res = PerformRequest(router, req)
    fmt.Println(res.Body.String())

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, fileContent, res.Body.Bytes())
}

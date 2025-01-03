package tests

import (
	// "bytes"
	"encoding/json"
	"fmt"

	"net/http"

	// "net/http/httptest"
	"testing"

	"github.com/blacac3/go-rest-api/internal/api"
	"github.com/blacac3/go-rest-api/internal/models"
	// "github.com/gin-gonic/gin"

	// "github.com/blacac3/go-rest-api/internal/util"
	"github.com/stretchr/testify/assert"
)

var user3 = models.User{
    FirstName: "Charlie", Surname: "Chaplin", 
    Username: "charliechaplin", Email: "charlie@example.com", 
    Password: "charlie123",
}

func TestJWT(t *testing.T) {
    SetupServer()
    router := api.InitRouter()

    var result map[string]interface{}
    registerUser(t, user3)
    resp := LoginUser(t, user3)
    response:=resp.Body.String()
        
    err := json.Unmarshal([]byte(response), &result)
    if err != nil{
        t.Errorf("Parse Error: %v", err)
    }

    token := result["token"]
    req, err := http.NewRequest("GET", baseUrl + "/healthz", nil)
    if err != nil {
        t.Fatalf("Failed to Create Request for health endpoint Test ---> %v", err)
    }
	req.Header.Set("Authorization", "Bearer "+ fmt.Sprint(token))
    
    res := PerformRequest(router, req)
    
    assert.Contains(t, res.Body.String(), "running", "Failed to validate JWT token")
    assert.Equal(t, http.StatusOK, res.Code, fmt.Sprintf("STATUS CODE:: Expected: 201, Got: %v", res.Code))
}



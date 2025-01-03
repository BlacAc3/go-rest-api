package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blacac3/go-rest-api/internal/api"
	"github.com/blacac3/go-rest-api/internal/models"
	// "github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Asignments
// ------------------------
var (
    User1 = models.User{
        FirstName: "Alice", Surname: "Goldman", 
        Username: "aliceinthelookingglass", Email: "alice@example.com", 
        Password: "alice123",
    }
    user2 = models.User{
        FirstName: "Bob", Surname: "Builder", 
        Username: "bobthebuilder", Email: "bob@example.com", 
        Password: "bob123",
    }
    user1Login = models.UserLogin{Email: "alice@example.com", Password: "alice123"}
    user2Login = models.UserLogin{Email: "bob@example.com", Password: "bob123"}

)

 



// Helper Functions 
// ----------------------------
func registerUser(t *testing.T, user models.User) *httptest.ResponseRecorder {
    jsonBody, err := json.Marshal(user)
    req, err := http.NewRequest("POST", baseUrl+"/register", bytes.NewBuffer(jsonBody))
    if err != nil {
        t.Fatalf("Failed to Create Request for Registration Test ---> %v", err)
    }
    router := api.InitRouter()
    res := PerformRequest(router, req)
    return res
}

func LoginUser(t *testing.T, user interface{}) *httptest.ResponseRecorder {
    jsonBody, err := json.Marshal(user)
    req, err := http.NewRequest("POST", baseUrl+"/login", bytes.NewBuffer(jsonBody))
    if err != nil {
        t.Fatalf("Failed to Create Request for Login Test ---> %v",err)
    }
    router := api.InitRouter()
    res := PerformRequest(router, req)
    return res
}




// Test Functions
// --------------------------------

//
func TestRegisterUser(t *testing.T) {
    SetupServer()
    res := registerUser(t, User1)
    assert.Equal(t, http.StatusCreated, res.Code, fmt.Sprintf("STATUS CODE:: Expected: 201, Got: %v", res.Code))
    assert.Contains(t, res.Body.String(), "aliceinthelookingglass", "Response Body does not contain username")
}



func TestSuccessfulLogin(t *testing.T) {
    SetupServer()
    res := LoginUser(t, User1)
    assert.Equal(t, http.StatusOK, res.Code, fmt.Sprintf("STATUS CODE:: Expected: %v, Got: %v", http.StatusOK, res.Code))
    assert.Contains(t, res.Body.String(), "token", "Response Body does not contain a token")
}


func TestFailedLogin(t *testing.T) {
    SetupServer()
    user2.Password = "bob1234"
    res := LoginUser(t, user2)
    assert.Equal(t, http.StatusNotFound, res.Code, fmt.Sprintf("STATUS CODE:: Expected: %v, Got: %v", http.StatusNotFound, res.Code))
    assert.Contains(t, res.Body.String(), "error", "Response Body does not contain an Error")
}

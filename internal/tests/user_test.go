package tests

import (
    "bytes"
    "net/http"
    "net/http/httptest"
    "testing"
    "encoding/json"

    "github.com/blacac3/go-rest-api/internal/models"
    "github.com/blacac3/go-rest-api/internal/api"
    "github.com/stretchr/testify/assert"
)

var user1 = models.User{FirstName: "Alice", Username: "aliceinthelookingglass", Email: "alice@example.com"}
var user2 =models.User{FirstName: "Bob", Username: "bobthebuilder", Email: "bob@example.com"}



func TestRegisterUser(t *testing.T){
    res := registerUser(t, user1)
    assert.Equal(t, http.StatusCreated, res.Code, "Expected Status Code 201 Received")
    assert.Contains(t, res.Body.String(), "aliceinthelookingglass", "Response Body does not contain username")
}

func registerUser(t *testing.T, user models.User) *httptest.ResponseRecorder{
    jsonBody, err := json.Marshal(user)
    req, err:= http.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
    if err != nil {
        t.Fatalf("Could not Create Request")
    }
    res := httptest.NewRecorder()
    api.HandleRegisteration(res, req)
    return res
}

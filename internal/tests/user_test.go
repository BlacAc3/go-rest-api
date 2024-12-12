package tests

import (
     "net/http"
     "net/http/httptest"
     "testing"
     "github.com/blacac3/go-rest-api/internal/models"
     "github.com/blacac3/go-rest-api/internal/api"
 )

var users = []models.User{
    {FirstName: "Alice", Username: "aliceinthelookingglass", Email: "alice@example.com"},
    {FirstName: "Bob", Username: "bobthebuilder", Email: "bob@example.com"},
}


func TestRegisterUser(t *testing.T){
    req, err:= http.NewRequest("POST", "/register", nil)
    if err != nil {
        t.Fatalf("Could not Create Request")
    }
    res := httptest.NewRecorder()
    api.HandleRegisteration(res, req)

    assert.Equal(t, http.StatusCreated, res.Code, "Expected Status Code 201 Received")
    assert.Contains(t, res.Body.String(), "aliceinthelookingglass", "Response Body does not contain username")

}

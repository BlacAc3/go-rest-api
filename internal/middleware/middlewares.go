package middleware

import (
	// "fmt"
	"log"
	"net/http"
	"strings"

	"github.com/blacac3/go-rest-api/internal/util"
	"github.com/gin-gonic/gin"
)


func Authentication() gin.HandlerFunc{
    return func(c *gin.Context){
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer "){
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
            c.Abort()
            log.Println("Authentication Failed: No Authorization Header")
            return
        }
        token := strings.TrimPrefix(authHeader, "Bearer ")
        _, err := util.VerifyJWT(token)
        if err!=nil{
            c.JSON(http.StatusUnauthorized, gin.H{"message":"Unauthorized"})
            c.Abort()
            log.Println("Authentication Failed: No Authorization Header")
            return
        }
        c.Next()
    }
}



 

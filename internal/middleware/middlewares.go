package middleware

import (
	// "fmt"
	"log"
	"net/http"

	"github.com/blacac3/go-rest-api/internal/util"
	"github.com/gin-gonic/gin"
)


func Authentication() gin.HandlerFunc{
    return func(c *gin.Context){
        token := util.GetJWT(c)
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



 

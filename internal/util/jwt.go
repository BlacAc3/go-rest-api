package util

import (

	// "errors"
    "net/http"
	"fmt"
	"strings"
	"log"
	"time"

	"github.com/blacac3/go-rest-api/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)


var secretKey = []byte("MySecretKey")

func GenerateJWT(user models.User) (string, error){
    claims := jwt.RegisteredClaims{
		Issuer:    "auth.example.com",
		Subject:   fmt.Sprint(user.ID),
		Audience:  []string{"go-rest-api"},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 5)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(secretKey)
    if err != nil{
        log.Printf("Error generating JWT: %v", err)
        return "", err
    }
    return tokenString, nil


}


func VerifyJWT(tokenString string) (string, error){
    // Verify Token
    token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{},func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

    if err != nil {
		return "", fmt.Errorf("failed to parse token: %v", err)
	}

    //Verify token Payload/Claims
    if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
        if !claims.VerifyExpiresAt(time.Now(), true) {
			return "", fmt.Errorf("token is expired")
		}
		if !claims.VerifyIssuer("auth.example.com", true) {
			return "", fmt.Errorf("invalid issuer")
		}
		if !claims.VerifyAudience("go-rest-api", true) {
			return "", fmt.Errorf("invalid audience")
		}

        return fmt.Sprint(claims.Subject), nil
	}
    return "", fmt.Errorf("invalid token")
} 


func GetJWT(c *gin.Context) string{
    authHeader := c.GetHeader("Authorization")
    if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer "){
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        c.Abort()
        log.Println("Authentication Failed: No Authorization Header")
        return ""
    }
    return strings.TrimPrefix(authHeader, "Bearer ")


}

package util

import (

	// "errors"
	"fmt"
	// "strings"
	"time"
    "log"

    "github.com/blacac3/go-rest-api/internal/models"
    "github.com/blacac3/go-rest-api/internal/database"
    "github.com/golang-jwt/jwt/v4"
)

var secretKey = []byte("MySecretKey")

func GenerateJWT(user models.User) (string, error){
    claims := jwt.RegisteredClaims{
		Issuer:    "auth.example.com",
		Subject:   string(user.ID),
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


func VerifyJWT(tokenString string) (interface{}, error){
    var user models.User 
    // Verify Token
    token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{},func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

    if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

    //Verify token Payload/Claims
    if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
        if !claims.VerifyExpiresAt(time.Now(), true) {
			return nil, fmt.Errorf("token is expired")
		}
		if !claims.VerifyIssuer("auth.example.com", true) {
			return nil, fmt.Errorf("invalid issuer")
		}
		if !claims.VerifyAudience("go-rest-api", true) {
			return nil, fmt.Errorf("invalid audience")
		}

		fmt.Println("Token is valid. Claims:", claims)
        result := database.DB.Where("id = ?", claims.Subject).First(&user)
		return result, nil
	}
    return nil, fmt.Errorf("invalid token")
} 

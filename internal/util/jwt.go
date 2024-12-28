package util

import (

	// "errors"
	"fmt"
	// "strings"
	"time"
    "log"

    "github.com/blacac3/go-rest-api/internal/models"
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
    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        if exp, ok := claims["exp"].(float64); ok {
			if time.Unix(int64(exp), 0).Before(time.Now()) {
				return nil, fmt.Errorf("token is expired")
			}
		}

    } else {
        return nil, fmt.Errorf("invalid token")
    }


} 

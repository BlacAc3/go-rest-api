package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	// "fmt"
	"strings"
	"time"
)

var secretKey = []byte("MySecretKey")

type Header struct {
    Alg string `json:"alg"`
    Typ string `json:"typ"`
}

type Payload struct {
    UserID int      `json:"user_id"`
    Role string     `json:"role"`
    Exp time.Time   `json:"exp"`
}




func EncodeBase64(data interface{}) (string,error){
    jsonData, err := json.Marshal(data)
    if err != nil{
        return "", err
    }
    return base64.RawURLEncoding.EncodeToString(jsonData), nil
}


func DecodeBase64(encoded string, out interface{}) error{
    data , err:= base64.RawURLEncoding.DecodeString(encoded)
    if err != nil {
        return err
    }
    return json.Unmarshal(data, out)
}



func CreateSignature(header64, payload64 string) string{
    h := hmac.New(sha256.New, secretKey)
    h.Write([]byte(header64 + "." + payload64))
    return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

func GenerateJWT(headerStruct Header, payloadStruct Payload) (string, error){
    header64, err := EncodeBase64(headerStruct)
    payload64, err := EncodeBase64(payloadStruct)
    if err != nil {
        return "", err
    }

    signature:= CreateSignature(header64, payload64)
    return header64 + "." + payload64 + "." + signature, nil
}


func VerifyJWT(token string) (interface{}, error){
    parts := strings.Split(token, ".")
    if len(parts) != 3 {
        return nil, errors.New("Invalid JWT format") 
    }

    header64 := parts[0]
    payload64 := parts[1]
    payloadData :=&Payload{}
    
    receivedSignature, err := base64.RawURLEncoding.DecodeString(parts[2])
    newSig, err := base64.RawURLEncoding.DecodeString(CreateSignature(header64, payload64))
    if err != nil{
        return nil, err
    }

    if hmac.Equal(newSig, receivedSignature) == false{
        return nil, errors.New("Invalid JWT signature")
    }
    return DecodeBase64(payload64, payloadData), nil
}
    

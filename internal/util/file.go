package util

import (
    "fmt"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

var encryptionKey string= "thisisasecretkey"
func EncryptFile(data []byte) (string, error){
    block, err := aes.NewCipher([]byte(encryptionKey))
    if err!=nil{
        return "", err
    }

    aesGCM, err := cipher.NewGCM(block)
    if err!=nil{
        return "", err
    }

    nonce := make([]byte, aesGCM.NonceSize())
    cipherText := aesGCM.Seal(nonce, nonce, data, nil)
    return base64.StdEncoding.EncodeToString(cipherText), nil
}



func DecryptFile(encryptedData string) ([]byte, error) {
    // Decode the Base64-encoded ciphertext
    cipherText, err := base64.StdEncoding.DecodeString(encryptedData)
    if err != nil {
        return nil, err
    }

    // Create a new AES cipher block
    block, err := aes.NewCipher([]byte(encryptionKey))
    if err != nil {
        return nil, fmt.Errorf("Error creating cipher block %v",err)
    }

    // Create a GCM instance
    aesGCM, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }

    // Extract the nonce from the ciphertext
    nonceSize := aesGCM.NonceSize()
    if len(cipherText) < nonceSize {
        return nil, fmt.Errorf("ciphertext too short")
    }
    nonce, cipherText := cipherText[:nonceSize], cipherText[nonceSize:]

    // Decrypt the data
    decryptedData, err := aesGCM.Open(nil, nonce, cipherText, nil)
    if err != nil {
        return nil, err
    }

    return decryptedData, nil
}


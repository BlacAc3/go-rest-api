package util

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

var encryptionKey string= "Ace the ultimate coder"
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

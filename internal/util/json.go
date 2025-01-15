package util

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RespondWithJson(c *gin.Context, statusCode int, payload interface{}){
    jsonData, err := json.Marshal(payload)
    if err != nil {
        http.Error(c.Writer, "Unable to encode JSON", http.StatusInternalServerError)   
        return
    }

    c.JSON(statusCode, gin.H{"payload": string(jsonData)})
}

// prettifyJSON takes an object and returns its prettified JSON string
func prettifyJSON(data interface{}) (string, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "  ") // Set indentation for prettifying
	err := encoder.Encode(data)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}

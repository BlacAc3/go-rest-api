package util

import (
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


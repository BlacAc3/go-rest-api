package util

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)
var validate = validator.New()


func ValidateRequest(c *gin.Context, model interface{}) error{
    decoder := json.NewDecoder(c.Request.Body)
    if err := decoder.Decode(model); err != nil {
        return fmt.Errorf("%v", err)
    }

    if err:= validate.Struct(model); err != nil{
        return fmt.Errorf("%v", err)
    }
    defer c.Request.Body.Close()

    return nil
}

package util

import (
    "fmt"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)
var validate = validator.New()


func ValidateRequest(r http.Request, model interface{}) error{
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(model); err != nil {
        return fmt.Errorf("%v", err)
    }

    if err:= validateStruct(model); err != nil{
        return fmt.Errorf("%v", err)
    }
    defer r.Body.Close()

    return nil
}

func validateStruct(s interface{}) error {
    return validate.Struct(s)

}

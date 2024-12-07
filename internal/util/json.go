package util
    
import(
    "net/http"
    "encoding/json"
)

func RespondWithJson(w http.ResponseWriter, statusCode int, payload interface{}){
    jsonData, err := json.Marshal(payload)
    if err != nil {
        http.Error(w, "Unable to encode JSON", http.StatusInternalServerError)   
        return
    }


    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    w.Write(jsonData)
}

package middleware

import (
	"log"
    "net/http"
)

type wrappedWriter struct {
    http.ResponseWriter
    statusCode int
}


func (w *wrappedWriter) WriteHeader(statusCode int) {
    w.ResponseWriter.WriteHeader(statusCode)
    w.statusCode = statusCode
}


func Logging(next http.Handler) http.Handler{
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
        wrapped := &wrappedWriter{
            ResponseWriter: w,
            statusCode: http.StatusOK,
        }
        next.ServeHTTP(wrapped, r)
        log.Println(wrapped.statusCode, r.Method, r.URL.Path)
        
    })
}

func Authentication(next http.Handler) http.Handler{
    return http.HandlerFunc (func(w http.ResponseWriter, r *http.Request){
        if token := r.Header.Get("Authorization"); token!="Bearer token"{
            http.Error(w, "{'message':'Unauthorized'}", http.StatusUnauthorized) 
            return
        }else{
            next.ServeHTTP(w, r)        
        }
        
    })
}



 

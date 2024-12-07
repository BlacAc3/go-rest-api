package api

import (
	"log"
	"net/http"
    "github.com/blacac3/go-crud-api/internal/middleware"

)


type APIServer struct{
    Port string

}

func NewAPIServer(p string) *APIServer{
    return &APIServer{Port: p}
}


func (s *APIServer) Serve(){
    router := http.NewServeMux()

    //Assign free Routes
    router.HandleFunc("POST /login", HandleLogin)
    router.HandleFunc("POST /register", HandleRegisteration)
    router.HandleFunc("/healthz", HandleHealthz)
    
    // Assign protected routes
    authRouter := http.NewServeMux()
    router.Handle("/", middleware.Authentication(authRouter))

    v1 := http.NewServeMux()
    v1.Handle("/v1/", http.StripPrefix("/v1", router))
    
    // Assign Middlewares
    applyMiddleware := middleware.CreateMiddlewareStack(
        middleware.Logging,
    )
    
    // Configure and start server
    server := http.Server{
        Addr: ":"+s.Port,
        Handler: applyMiddleware(v1),
    }
    log.Fatal(server.ListenAndServe())
}

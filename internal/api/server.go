package api

import (
	"log"
	"net/http"

	"github.com/blacac3/go-rest-api/internal/middleware"
	"github.com/gin-gonic/gin"
)


type APIServer struct{
    Port string

}

func NewAPIServer(p string) *APIServer{
    return &APIServer{Port: p}
}

func InitRouter() *gin.Engine{
    router:= gin.Default()
    router.POST("/login", HandleLogin)
    router.POST("/register", HandleRegisteration)

    authGroup := router.Group("/")
    authGroup.Use(middleware.Authentication())
    {
        authGroup.GET("/healthz", HandleHealthz)
    }

    // Assign Middlewares
    return router

}


func (s *APIServer) Serve(){
    router := InitRouter()
    // Configure and start server
    server := http.Server{
        Addr: ":"+s.Port,
        Handler: router,
    }
    log.Fatal(server.ListenAndServe())
}

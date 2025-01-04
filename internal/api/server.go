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
    router.POST("auth/login", HandleLogin)
    router.POST("auth/register", HandleRegisteration)

    authGroup := router.Group("/")
    authGroup.Use(middleware.Authentication())
    {
        authGroup.GET("/healthz", HandleHealthz)
        authGroup.POST("/files/upload", HandleFileUpload)
        authGroup.GET("/files/:fileID", HandleFileDownload)
    }

    // Assign Middlewares
    return router

}


func (s *APIServer) Serve(){
    gin.SetMode(gin.ReleaseMode)
    router := InitRouter()
    // Configure and start server
    server := http.Server{
        Addr: ":"+s.Port,
        Handler: router,
    }
    log.Fatal(server.ListenAndServe())
}

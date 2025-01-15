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



func InitRouter(router *gin.Engine) {
    router.POST("auth/login", HandleLogin)
    router.POST("auth/register", HandleRegisteration)

    authGroup := router.Group("/")
    authGroup.Use(middleware.Authentication())
    {
        authGroup.GET("/healthz", HandleHealthz)
        authGroup.POST("/files/upload", HandleFileUpload)
        authGroup.GET("/files/:fileID", HandleFileDownload)
        authGroup.GET("/files", HandleFilesRetrieval)
    }

}


func (s *APIServer) Serve(){
    router := gin.Default()
    InitRouter(router)
    // Configure and start server
    server := http.Server{
        Addr: ":"+s.Port,
        Handler: router,
    }
    log.Fatal(server.ListenAndServe())
}

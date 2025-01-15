package tests

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"github.com/blacac3/go-rest-api/internal/api"
	"github.com/blacac3/go-rest-api/internal/database"
	"github.com/gin-gonic/gin"
)
var (
    port string = "8100"
    baseUrl string = fmt.Sprintf("http://localhost:%s", port)
    serverStatus = ""
)




// Start Test Server
func SetupServer() *http.Server {
    workingDir, err := database.FindProjectRoot("go.mod")
    testDBPath := fmt.Sprintf("%v/%v", workingDir, database.TestDBName)
    if err != nil {
        log.Printf("Could not find project root due to the following error: %v", err)
        return nil
    }
    _, err = os.Stat(testDBPath)
    if err == nil {
        err = os.Remove(fmt.Sprintf("%v/%v", workingDir, database.TestDBName))
        if err != nil {
            log.Printf("Could not remove test database due to the following error: %v", err)
            return nil
        }
    }


    gin.SetMode(gin.TestMode)
    database.Test_Mode = true
    if serverStatus == "" {
        go func() {
            router := GetRouter()
            if err := router.Run(":8000"); err != nil {
                panic(err)
            }
        }()
        time.Sleep(1 * time.Second)
        serverStatus = "Live"
    }
    return nil
}

func GetRouter() *gin.Engine{
    router:=gin.New()
    router.Use(gin.Recovery())
    api.InitRouter(router)
    return router

}


func PerformRequest(router *gin.Engine, req *http.Request) *httptest.ResponseRecorder {
    gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

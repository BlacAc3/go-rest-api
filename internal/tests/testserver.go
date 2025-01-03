package tests

import (
	"fmt"
	"net/http"
    "net/http/httptest"
	"time"

	"github.com/blacac3/go-rest-api/internal/database"
    "github.com/blacac3/go-rest-api/internal/api"
	"github.com/gin-gonic/gin"

)
var (
    port string = "8100"
    baseUrl string = fmt.Sprintf("http://localhost:%s", port)
    serverStatus = ""
)




// Start Test Server
func SetupServer() *http.Server {
    gin.SetMode(gin.ReleaseMode)
    database.Test_Mode = true
    if serverStatus == "" {
        server := api.NewAPIServer(port)
        go func() {
            server.Serve()
        }()
        time.Sleep(1 * time.Second)
        serverStatus = "Live"
    }
    return nil
}




func PerformRequest(router *gin.Engine, req *http.Request) *httptest.ResponseRecorder {
    gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

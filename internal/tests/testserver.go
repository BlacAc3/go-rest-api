package tests

import (
	"fmt"
	"net/http"
	"time"

	"github.com/blacac3/go-rest-api/internal/api"
	"github.com/blacac3/go-rest-api/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)
var (
    port string = "8100"
    baseUrl string = fmt.Sprintf("http://localhost:%s", port)
    serverStatus = ""
    db *gorm.DB 
)



// Start Test Server
func SetupServer() *http.Server {
    if serverStatus == "" {
        db, _ = setupTestDB()
        server := api.NewAPIServer(port)
        go func() {
            server.Serve()
        }()
        time.Sleep(1 * time.Second)
        serverStatus = "Live"
    }
    return nil
}



func setupTestDB() (*gorm.DB, error) {
	// Use SQLite in-memory database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Silent logger for testing
	})
	if err != nil {
		return nil, err
	}

	// Run migrations
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		return nil, err
	}
    api.ChangeDB(db)

	return db, nil
}

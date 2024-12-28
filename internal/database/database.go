package database

import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
	"gorm.io/gorm/logger"
    "github.com/blacac3/go-rest-api/internal/models"
	"log"
)

var DB *gorm.DB

func ConnectToDB(){
    var err error
    dbString := "postgres://postgres:blacac3@localhost:5432/goapidb?sslmode=disable"
    DB, err = gorm.Open(postgres.Open(dbString), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Silent), // Suppress logs
    })


    if err != nil {
        log.Println("Failed to Connect to Database:")
        log.Fatal(err)
    }else{
        log.Println("Connected to Database Successfully✅.")
    }
    
    //Start Transaction
    tx := DB.Begin()
    log.Println("Starting Transaction for new migration.")

    if err = tx.AutoMigrate(&models.User{}); err!=nil {
        log.Printf("⚠⚠ Failed to make Migrations to Database: %v",err)
        tx.Rollback()
        log.Fatal("‼‼ Rolling back Transaction")
    }else{
        log.Println("Migrations ran Successfully✅.")
    }

    // Commit the transaction
    if err := tx.Commit().Error; err != nil {
        log.Fatalf("Failed to commit transaction: %v", err)
        return
    }
    log.Println("Transaction Successful✅.")

    log.Println("DB Connection and Migration Successful✅.")
}


// GetDB returns the global database instance
func GetDB() *gorm.DB {
	if DB == nil {
		log.Fatal("Database is not connected. Call db.Connect() first.")
	}
	return DB
}

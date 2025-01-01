package database

import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
	"gorm.io/gorm/logger"
    "github.com/blacac3/go-rest-api/internal/models"
	"log"
)

var sqlDB *gorm.DB

func ChangeDB(db *gorm.DB){
    sqlDB = db
}

func ConnectToDB(){
    var err error
    dbString := "postgres://postgres:blacac3@localhost:5432/goapidb?sslmode=disable"
    sqlDB, err = gorm.Open(postgres.Open(dbString), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Silent), // Suppress logs
    })


    if err != nil {
        log.Fatalf("Failed to Connect to Database: %v", err)
    }else{
        log.Println("Connected to Database Successfully✅.")
    }
    
    //Start Transaction
    tx := sqlDB.Begin()
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

    log.Println("sqlDB Connection and Migration Successful✅.")
}



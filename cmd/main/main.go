package main

import (
	// net/http"
	"log"

	"github.com/blacac3/go-rest-api/internal/api"
    // "github.com/blacac3/go-rest-api/internal/database"

)

func main() {
	var port string = "8000"
	server := api.NewAPIServer(port)

    log.Println()
    log.Println("📦📦Server started on port", port)
	server.Serve()

}

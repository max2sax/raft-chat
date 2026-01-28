package main

import (
	"fmt"
	"net/http"

	"github.com/max2sax/raft-chat/api"
	"github.com/max2sax/raft-chat/storage"
)

func main() {
	// Initialize storage
	store := storage.NewStorage()

	// Initialize HTTP server
	server := &http.Server{
		Addr: ":8080",
	}

	// Initialize API
	chatAPI := api.NewAPI(store, server).
		RegisterRoutes()

	// Start server
	fmt.Println("Server starting on :8080")
	if err := chatAPI.Start(); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

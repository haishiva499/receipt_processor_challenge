package main

import (
	"fmt"
	"log"
	"net/http"
	"receipt-processor/internal/receipt"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize the router
	r := mux.NewRouter()

	// Define the routes
	r.HandleFunc("/receipts/process", receipt.ProcessReceipt).Methods("POST")
	r.HandleFunc("/receipts/{id}/points", receipt.GetPoints).Methods("GET")

	// Start the server
	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

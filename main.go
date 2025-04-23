package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/warrenu/receipt-processor/cache"
	"github.com/warrenu/receipt-processor/handlers"
)

func main() {
	// Initialize the cache: Can support 10,000 unique records...
	store := cache.NewStore[int](10000)

	// Initialize the handler (which internally uses the service logic)
	handler := handlers.NewHandler(store)

	// Set up the router and routes
	router := mux.NewRouter()

	// Middleware to see requests that come in:
	// On Localhost, this most likely would be: `[::1]:<port>` a IPv6 loopback address
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("%s %s %s\n", r.Method, r.RequestURI, r.RemoteAddr)
			next.ServeHTTP(w, r)
		})
	})

	// Endpoints:
	router.HandleFunc("/receipts/process", handler.ProcessReceipt).Methods("POST")
	router.HandleFunc("/receipts/{id}/points", handler.GetPoints).Methods("GET")

	// Start the server
	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}

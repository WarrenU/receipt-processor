package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/warrenu/receipt-processor/cache"
	"github.com/warrenu/receipt-processor/models"
	"github.com/warrenu/receipt-processor/service"
)

// ServiceHandler handles HTTP requests for receipts
// backed by a generic LRU cache of point values.
type ServiceHandler struct {
	Store *cache.Store[int]
}

// NewHandler creates a new ServiceHandler with the given int cache.
func NewHandler(store *cache.Store[int]) *ServiceHandler {
	return &ServiceHandler{Store: store}
}

// ProcessReceipt parses a receipt, uses service.ProcessReceipt to compute
// an ID and points, caches the points, and returns the ID.
func (h *ServiceHandler) ProcessReceipt(w http.ResponseWriter, r *http.Request) {
	var receipt models.Receipt
	if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	id, points, err := service.ProcessReceipt(receipt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.Store.Set(id, points)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

// GetPoints retrieves cached points for the given receipt ID.
func (h *ServiceHandler) GetPoints(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	points, ok := h.Store.Get(id)
	if !ok {
		http.Error(w, "ID not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"points": points})
}

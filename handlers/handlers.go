package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/warrenu/receipt-processor/cache"
	"github.com/warrenu/receipt-processor/models"
	"github.com/warrenu/receipt-processor/service"
)

type ServiceHandler struct {
	Store *cache.Store[int]
}

func NewHandler(store *cache.Store[int]) *ServiceHandler {
	return &ServiceHandler{Store: store}
}

func (h *ServiceHandler) ProcessReceipt(w http.ResponseWriter, r *http.Request) {
	var receipt models.Receipt
	if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}

	receipt.Sanitize()

	validationErrors := ValidateStruct(receipt)
	if validationErrors != nil {
		respondJSON(w, http.StatusBadRequest, validationErrors)
		return
	}

	id, points, err := service.ProcessReceipt(receipt)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to process receipt: "+err.Error())
		return
	}

	h.Store.Set(id, points)

	respondJSON(w, http.StatusOK, map[string]string{"id": id})
}

func (h *ServiceHandler) GetPoints(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := strings.TrimSpace(vars["id"])

	points, ok := h.Store.Get(id)
	if !ok {
		respondError(w, http.StatusNotFound, "ID not found")
		return
	}

	respondJSON(w, http.StatusOK, map[string]int{"points": points})
}

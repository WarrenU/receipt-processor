package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/warrenu/receipt-processor/cache"
	"github.com/warrenu/receipt-processor/handlers" // Reference handlers instead of service
	"github.com/warrenu/receipt-processor/models"
)

func TestGetPoints(t *testing.T) {
	// 1) Init shared cache + handler
	store := cache.NewStore[int](100)
	handler := handlers.NewHandler(store)

	// 2) Build a Gorilla mux router (same as main.go)
	router := mux.NewRouter()
	router.HandleFunc("/receipts/process", handler.ProcessReceipt).Methods("POST")
	router.HandleFunc("/receipts/{id}/points", handler.GetPoints).Methods("GET")

	// 3) Start a test server
	server := httptest.NewServer(router)
	defer server.Close()

	// 4) POST a receipt to get an ID
	receipt := models.Receipt{
		Retailer:     "Test Store",
		PurchaseDate: "2025-03-20",
		PurchaseTime: "15:00",
		Items: []models.Item{
			{ShortDescription: "Item A", Price: "2.25"},
			{ShortDescription: "Item B", Price: "3.75"},
		},
		Total: "6.00",
	}
	body, _ := json.Marshal(receipt)
	postResp, err := http.Post(server.URL+"/receipts/process",
		"application/json", bytes.NewBuffer(body))
	assert.NoError(t, err)
	defer postResp.Body.Close()

	var idResp map[string]string
	json.NewDecoder(postResp.Body).Decode(&idResp)
	id := idResp["id"]
	assert.NotEmpty(t, id)

	// 5) GET the points back
	getResp, err := http.Get(server.URL + "/receipts/" + id + "/points")
	assert.NoError(t, err)
	defer getResp.Body.Close()

	assert.Equal(t, http.StatusOK, getResp.StatusCode)
	var pts map[string]int
	json.NewDecoder(getResp.Body).Decode(&pts)
	assert.Greater(t, pts["points"], 0)
	t.Log("API Integration tests passed succesfully")
}

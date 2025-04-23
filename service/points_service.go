package service

import (
	"log"

	"github.com/warrenu/receipt-processor/models"
)

// ProcessReceipt generates a deterministic ID for the receipt
// and computes its points, failing if UUID generation fails.
func ProcessReceipt(receipt models.Receipt) (string, int, error) {
	// Generate the deterministic UUID
	id, err := GetDeterministicUUID(receipt)
	if err != nil {
		log.Printf("Error generating deterministic UUID: %v", err)
		return "", 0, err // Return empty ID and 0 points if UUID generation fails
	}

	// Calculate points
	points, _ := CalculatePoints(receipt)

	return id, points, nil
}

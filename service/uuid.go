package service

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/warrenu/receipt-processor/models"
)

// GetDeterministicUUID generates a deterministic UUID based on the receipt's contents.
// Returns an error if marshalling the receipt fails.
func GetDeterministicUUID(receipt models.Receipt) (string, error) {
	// Marshal the receipt into bytes
	data, err := json.Marshal(receipt)
	if err != nil {
		log.Printf("error marshaling receipt for UUID: %v", err)
		return "", errors.New("failed to marshal receipt for deterministic UUID")
	}

	// Generate UUID using SHA1 (v5) with a custom namespace
	namespace := uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	deterministicUUID := uuid.NewSHA1(namespace, data)

	// Return the UUID as a string
	return deterministicUUID.String(), nil
}

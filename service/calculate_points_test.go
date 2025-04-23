package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/warrenu/receipt-processor/models"
	"github.com/warrenu/receipt-processor/service"
)

func TestCalculatePoints(t *testing.T) {
	tests := []struct {
		name           string
		receipt        models.Receipt
		expectedPoints int
	}{
		// Original example tests
		{
			name: "Target Example",
			receipt: models.Receipt{
				Retailer:     "Target",
				PurchaseDate: "2022-01-01",
				PurchaseTime: "13:01",
				Items: []models.Item{
					{ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
					{ShortDescription: "Emils Cheese Pizza", Price: "12.25"},
					{ShortDescription: "Knorr Creamy Chicken", Price: "1.26"},
					{ShortDescription: "Doritos Nacho Cheese", Price: "3.35"},
					{ShortDescription: "Klarbrunn 12-PK 12 FL OZ", Price: "12.00"},
				},
				Total: "35.35",
			},
			expectedPoints: 28,
		},
		{
			name: "M&M Corner Market Example",
			receipt: models.Receipt{
				Retailer:     "M&M Corner Market",
				PurchaseDate: "2022-03-20",
				PurchaseTime: "14:33",
				Items: []models.Item{
					{ShortDescription: "Gatorade", Price: "2.25"},
					{ShortDescription: "Gatorade", Price: "2.25"},
					{ShortDescription: "Gatorade", Price: "2.25"},
					{ShortDescription: "Gatorade", Price: "2.25"},
				},
				Total: "9.00",
			},
			expectedPoints: 109,
		},

		// Additional edge and mixed cases
		{
			name: "Empty Receipt",
			receipt: models.Receipt{
				Retailer:     "",
				PurchaseDate: "2022-01-02",
				PurchaseTime: "10:00",
				Items:        []models.Item{},
				Total:        "0.00",
			},
			expectedPoints: 75,
		},
		{
			name: "Basic Retailer and Quarter",
			receipt: models.Receipt{
				Retailer:     "A1",
				PurchaseDate: "2022-01-02",
				PurchaseTime: "11:00",
				Items:        []models.Item{},
				Total:        "0.50",
			},
			expectedPoints: 27,
		},
		{
			name: "Item Description Bonus and Odd Day",
			receipt: models.Receipt{
				Retailer:     "",
				PurchaseDate: "2022-01-03",
				PurchaseTime: "09:00",
				Items: []models.Item{
					{ShortDescription: "ABCDEF", Price: "5.00"}, // length6*0.2=1.2 ceil2? actually 5*0.2=1.0 so ceil1
				},
				Total: "0.00",
			},
			expectedPoints: 82,
		},
		{
			name: "Even Day Afternoon",
			receipt: models.Receipt{
				Retailer:     "",
				PurchaseDate: "2022-01-04",
				PurchaseTime: "14:00",
				Items:        []models.Item{},
				Total:        "1.00",
			},
			expectedPoints: 85,
		},
		{
			name: "Three Items Pair Bonus and Quarter",
			receipt: models.Receipt{
				Retailer:     "",
				PurchaseDate: "2022-01-05",
				PurchaseTime: "12:00",
				Items: []models.Item{
					{ShortDescription: "Item1", Price: "2.00"},
					{ShortDescription: "Item2", Price: "3.00"},
					{ShortDescription: "Item3", Price: "1.00"},
				},
				Total: "3.00",
			},
			expectedPoints: 86,
		},
		{
			name: "No Bonus Scenario",
			receipt: models.Receipt{
				Retailer:     "",
				PurchaseDate: "2022-01-06",
				PurchaseTime: "08:00",
				Items:        []models.Item{},
				Total:        "4.10",
			},
			expectedPoints: 0,
		},
		{
			name: "Special Char Retailer",
			receipt: models.Receipt{
				Retailer:     "@#$",
				PurchaseDate: "2022-01-02",
				PurchaseTime: "10:00",
				Items:        []models.Item{},
				Total:        "0.00",
			},
			expectedPoints: 75,
		},
		{
			name: "Complex Mix",
			receipt: models.Receipt{
				Retailer:     "GoLang123!", // 9 points
				PurchaseDate: "2022-01-07", // 6 points
				PurchaseTime: "14:30",      // 10 points, 2:30PM purchase
				Items: []models.Item{
					{ShortDescription: "AltItem", Price: "6.00"},
					{ShortDescription: "BBB", Price: "2.50"},    // 2.5 * 0.2 = 0.05, rounded to nearest int = 1 point
					{ShortDescription: "No", Price: "1.00"},     // 1 * 0.2 = 0.2, rounded to nearest int = 1 point
					{ShortDescription: "YesYes", Price: "3.00"}, // 10 points for every two items on receipt
				},
				Total: "12.75", // 25 points, multiple of .25
			},
			expectedPoints: 62, // 9 + 6 + 10 + 1 + 1 + 10 + 25
		},
		{
			name: "Max Description Multiple",
			receipt: models.Receipt{
				Retailer:     "",
				PurchaseDate: "2022-01-02",
				PurchaseTime: "10:00",
				Items: []models.Item{
					{ShortDescription: "AAA", Price: "1.00"},
					{ShortDescription: "BBBBBB", Price: "2.00"},
				},
				Total: "0.00",
			},
			expectedPoints: 82,
		},
		{
			name: "Minimal Amount",
			receipt: models.Receipt{
				Retailer:     "",
				PurchaseDate: "2022-01-02",
				PurchaseTime: "10:00",
				Items:        []models.Item{},
				Total:        "0.01",
			},
			expectedPoints: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, breakdown := service.CalculatePoints(tt.receipt)
			t.Log(breakdown) // Log description of test case
			assert.Equal(t, tt.expectedPoints, got)
		})
	}
}

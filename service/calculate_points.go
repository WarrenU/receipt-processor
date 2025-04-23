package service

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/warrenu/receipt-processor/models"
)

// CalculatePoints applies all receipt-scoring rules and returns the total and a breakdown.
// breakdown as of now is just for test case output, to verify and validate how the receipt points
// are being calculated.
func CalculatePoints(receipt models.Receipt) (int, string) {
	points := 0
	breakdown := []string{""}

	rp := calculateRetailerPoints(receipt.Retailer)
	points += rp
	breakdown = append(breakdown, fmt.Sprintf("%d pts - Retailer Name", rp))

	rdp := calculateTotalRoundDollarPoints(receipt.Total)
	points += rdp
	breakdown = append(breakdown, fmt.Sprintf("%d pts - Round Dollar Total", rdp))

	qp := calculateTotalMultipleOfQuarterPoints(receipt.Total)
	points += qp
	breakdown = append(breakdown, fmt.Sprintf("%d pts - Multiple of $0.25", qp))

	icp := calculateItemCountPoints(receipt.Items)
	points += icp
	breakdown = append(breakdown, fmt.Sprintf("%d pts - Item Count every 2 items", icp))

	idp := calculateItemDescriptionPoints(receipt.Items)
	points += idp
	breakdown = append(breakdown, fmt.Sprintf("%d pts - Item Descriptions mult of 3", idp))

	odp := calculateOddDayPoints(receipt.PurchaseDate)
	points += odp
	breakdown = append(breakdown, fmt.Sprintf("%d pts - Odd Purchase Date", odp))

	ptp := calculatePurchaseTimePoints(receipt.PurchaseTime)
	points += ptp
	breakdown = append(breakdown, fmt.Sprintf("%d pts - Purchase Time", ptp))

	return points, strings.Join(breakdown, "\n")
}

// 1 point per alphanumeric character in the retailer name.
func calculateRetailerPoints(retailer string) int {
	re := regexp.MustCompile(`[A-Za-z0-9]`)
	return len(re.FindAllString(retailer, -1))
}

// parseCents converts a dollar string (e.g., "9.15") to its value in cents as an integer (915).
// This avoids floating-point errors by working with whole numbers.
func parseCents(amountStr string) (int, error) {
	parts := strings.SplitN(amountStr, ".", 2)
	dollars, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, err
	}
	cents := 0
	if len(parts) == 2 {
		frac := parts[1]
		if len(frac) == 1 {
			frac += "0"
		} else if len(frac) > 2 {
			frac = frac[:2]
		}
		cents, err = strconv.Atoi(frac)
		if err != nil {
			return 0, err
		}
	}
	return dollars*100 + cents, nil
}

// 50 points if total is a round dollar amount (no cents).
func calculateTotalRoundDollarPoints(totalStr string) int {
	cents, err := parseCents(totalStr)
	if err != nil {
		return 0
	}
	if cents%100 == 0 {
		return 50
	}
	return 0
}

// 25 points if total is a multiple of $0.25.
func calculateTotalMultipleOfQuarterPoints(totalStr string) int {
	cents, err := parseCents(totalStr)
	if err != nil {
		return 0
	}
	// 25 cents == a quarter
	if cents%25 == 0 {
		return 25
	}
	return 0
}

// 5 points for every two items on the receipt.
func calculateItemCountPoints(items []models.Item) int {
	return (len(items) / 2) * 5
}

// If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer.
func calculateItemDescriptionPoints(items []models.Item) int {
	points := 0
	for _, item := range items {
		desc := strings.TrimSpace(item.ShortDescription)
		if len(desc)%3 == 0 {
			price, err := strconv.ParseFloat(item.Price, 64)
			if err != nil {
				continue
			}
			points += int(math.Ceil(price * 0.2))
		}
	}
	return points
}

// 6 points if the purchase date’s day is odd.
func calculateOddDayPoints(dateStr string) int {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return 0
	}
	if t.Day()%2 == 1 {
		return 6
	}
	return 0
}

// 10 points if purchase time is after 2:00pm and before 4:00pm.
func calculatePurchaseTimePoints(timeStr string) int {
	t, err := time.Parse("15:04", timeStr)
	if err != nil {
		return 0
	}
	if t.Hour() == 14 {
		return 10
	}
	return 0
}

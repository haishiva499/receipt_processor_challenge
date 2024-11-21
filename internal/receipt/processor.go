package receipt

import (
	"encoding/json"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Item represents a product in the receipt.
type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

// Receipt represents a full receipt structure.
type Receipt struct {
	ID           string `json:"id,omitempty"`
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []Item `json:"items"`
	Total        string `json:"total"`
}

// Store for holding receipts in memory (in-memory storage for simplicity).
var receiptStore = make(map[string]Receipt)

// ProcessReceipt handles the submission of receipts and returns a unique ID.
func ProcessReceipt(w http.ResponseWriter, r *http.Request) {
	var receipt Receipt
	if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Generate a unique ID for the receipt
	receipt.ID = uuid.New().String()

	// Store the receipt in the in-memory store
	receiptStore[receipt.ID] = receipt

	// Return the receipt ID as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"id": receipt.ID})
}

// GetPoints calculates and returns the points awarded for a specific receipt.
func GetPoints(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Retrieve the receipt from the store
	receipt, exists := receiptStore[id]
	if !exists {
		http.Error(w, "Receipt not found", http.StatusNotFound)
		return
	}

	// Calculate the points
	points := calculatePoints(receipt)

	// Return the points as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int64{"points": points})
}

// calculatePoints calculates the points based on various rules.
func calculatePoints(receipt Receipt) int64 {
	points := int64(0)

	// Rule 1: Points for retailer name length (alphanumeric characters only)
	points += int64(len(getAlphanumericString(receipt.Retailer)))

	// Rule 2: 50 points if total is a round dollar amount with no cents
	if isRoundDollar(receipt.Total) {
		points += 50
	}

	// Rule 3: 25 points if total is a multiple of 0.25
	if isMultipleOfQuarter(receipt.Total) {
		points += 25
	}

	// Rule 4: Points for number of items (5 points for every 2 items)
	points += int64((len(receipt.Items) / 2) * 5)

	// Rule 5: Points based on item description length (multiple of 3)
	for _, item := range receipt.Items {
		if len(strings.TrimSpace(item.ShortDescription))%3 == 0 {
			price, _ := parsePrice(item.Price)
			points += int64(math.Ceil(price * 0.2)) // Rounded up price * 0.2
		}
	}

	// Rule 6: Points if purchase date day is odd
	if isOddDay(receipt.PurchaseDate) {
		points += 6
	}

	// Rule 7: Points if purchase time is between 2:00pm and 4:00pm
	if isBetween2And4PM(receipt.PurchaseTime) {
		points += 10
	}

	return points
}

// getAlphanumericString removes non-alphanumeric characters and returns a string of alphanumeric characters only.
func getAlphanumericString(s string) string {
	re := regexp.MustCompile("[^a-zA-Z0-9]")
	return re.ReplaceAllString(s, "")
}

// isRoundDollar checks if the total is a round dollar amount (e.g., "10.00")
func isRoundDollar(total string) bool {
	// Check if total is a round dollar amount (e.g., "10.00")
	return strings.HasSuffix(total, ".00")
}

// isMultipleOfQuarter checks if the total is a multiple of 0.25.
func isMultipleOfQuarter(total string) bool {
	// Convert the string total to float and check if divisible by 0.25
	price, err := parsePrice(total)
	if err != nil {
		return false
	}
	return math.Mod(price, 0.25) == 0
}

// parsePrice parses a price string into a float64 value.
func parsePrice(price string) (float64, error) {
	return strconv.ParseFloat(price, 64)
}

// isOddDay checks if the day of the purchase date is odd.
func isOddDay(date string) bool {
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return false
	}
	return parsedDate.Day()%2 != 0
}

// isBetween2And4PM checks if the purchase time is between 2:00 PM and 4:00 PM.
func isBetween2And4PM(timeStr string) bool {
	parsedTime, err := time.Parse("15:04", timeStr)
	if err != nil {
		return false
	}
	return parsedTime.Hour() >= 14 && parsedTime.Hour() < 16
}

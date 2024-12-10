package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
)

const RecNameSpace = "123e4567-e89b-12d3-a456-426614174000"

// ParseReceiptJson takes the json string of a receipt and unmarshals it into a struct.
func ParseReceiptJson(body io.Reader) (*Receipt, error) {
	var r Receipt
	if err := json.NewDecoder(body).Decode(&r); err != nil {
		return nil, err
	}
	return &r, nil
}

// GenerateUUIDForReceipt creates a deterministic UUID for a given receipt.
// Like a hash, but conforms to UUID
// NOTE: This is assuming that two identical receipts should be considered as duplicates.
func GenerateUUIDForReceipt(receipt *Receipt) string {
	var data []byte
	data = append(data, []byte(receipt.Retailer)...)
	data = append(data, []byte(receipt.PurchaseDate.Format(time.RFC3339))...)
	data = append(data, []byte(receipt.PurchaseTime.Format(time.RFC3339))...)
	data = append(data, []byte(fmt.Sprintf("%f", receipt.Total))...)
	for _, item := range receipt.Items {
		data = append(data, []byte(item.ShortDescription)...)
		data = append(data, []byte(fmt.Sprintf("%f", item.Price))...)
	}

	namespace := uuid.MustParse(RecNameSpace)
	receiptUUID := uuid.NewHash(sha1.New(), namespace, data, 5)
	return receiptUUID.String()
}

func CalculateReceiptPoints(r *Receipt) int {
	r.Points = calculateRetailerPoints(r.Retailer) +
		calculateTotalPoints(r.Total) +
		calculateItemPoints(r.Items) +
		calculateDatePoints(r.PurchaseDate) +
		calculateTimePoints(r.PurchaseTime)

	return r.Points
}

// One point for every alphanumeric char in retailer name
func calculateRetailerPoints(retailer string) int {
	points := 0
	for _, ch := range retailer {
		if isAlphanumeric(ch) {
			points++
		}
	}
	return points
}

// 50 points if the total is a round dollar amount with no cents.
// 25 points if the total is a multiple of 0.25
func calculateTotalPoints(total float64) int {
	points := 0
	if math.Mod(total, 1.0) == 0 {
		points += 50
	}
	if math.Mod(total, 0.25) == 0 {
		points += 25
	}
	return points
}

// 5 points for every two items on the receipt
// If the trimmed length of the item description is a multiple of 3,
// multiply the price by 0.2 and round up to the nearest integer.
// The result is the number of points earned.
func calculateItemPoints(items []Item) int {
	points := (len(items) / 2) * 5

	for _, item := range items {
		descriptionLength := len(strings.TrimSpace(item.ShortDescription))
		if descriptionLength%3 == 0 {
			points += int(math.Ceil(item.Price * 0.2))
		}
	}

	return points
}

// 6 points if the day in the purchase date is odd
func calculateDatePoints(purchaseDate DateTime) int {
	if purchaseDate.Day()%2 != 0 {
		return 6
	}
	return 0
}

// 10 points if the time of purchase is after 2pm and before 4pm
func calculateTimePoints(purchaseTime DateTime) int {
	hour := purchaseTime.Hour()
	if hour == 14 || (hour == 15 && purchaseTime.Minute() == 0) {
		return 10
	}
	return 0
}

func isAlphanumeric(ch rune) bool {
	return (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9')
}

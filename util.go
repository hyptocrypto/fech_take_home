package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	RecNameSpace = "123e4567-e89b-12d3-a456-426614174000"
	dateFormat   = "2006-01-02"
	timeFormat   = "15:04"
)

// ParseReceiptJson takes the json string of a receipt and unmarshals it into a struct.
// This must handle the non-standard time/date format
func ParseReceiptJson(body io.Reader) (*Receipt, error) {
	var temp StringReceipt
	if err := json.NewDecoder(body).Decode(&temp); err != nil {
		return nil, err
	}

	parsedDate, err := time.Parse(dateFormat, temp.PurchaseDate)
	if err != nil {
		return nil, fmt.Errorf("error parsing date: %w", err)
	}

	parsedTime, err := time.Parse(timeFormat, temp.PurchaseTime)
	if err != nil {
		return nil, fmt.Errorf("error parsing time: %w", err)
	}

	total, err := strconv.ParseFloat(temp.Total, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing total: %w", err)
	}

	receipt := &Receipt{
		Retailer:     temp.Retailer,
		PurchaseDate: parsedDate,
		PurchaseTime: parsedTime,
		Items:        temp.Items,
		Total:        total,
		Points:       -1,
	}

	return receipt, nil
}

// GenerateUUIDForReceipt creates a deterministic UUID for a given receipt.
// Like a hash, but conforms to UUID
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

// One point for every alphanumeric character in the retailer name.
// 50 points if the total is a round dollar amount with no cents.
// 25 points if the total is a multiple of 0.25.
// 5 points for every two items on the receipt.
// If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer. The result is the number of points earned.
// 6 points if the day in the purchase date is odd.
// 10 points if the time of purchase is after 2:00pm and before 4:00pm.
// TODO: This is pretty ugly. Re-address for code quality
func CalculateReceiptPoints(r *Receipt) int {
	points := 0

	for _, ch := range r.Retailer {
		if isAlphanumeric(ch) {
			points++
		}
	}

	if math.Mod(r.Total, 1.0) == 0 {
		points += 50
	}
	if math.Mod(r.Total, 0.25) == 0 {
		points += 25
	}

	points += (len(r.Items) / 2) * 5

	for _, item := range r.Items {
		descriptionLength := len(strings.TrimSpace(item.ShortDescription))
		if descriptionLength%3 == 0 {
			points += int(math.Ceil(item.Price * 0.2))
		}
	}

	if r.PurchaseDate.Day()%2 != 0 {
		points += 6
	}

	if r.PurchaseTime.Hour() == 14 || (r.PurchaseTime.Hour() == 15 && r.PurchaseTime.Minute() == 0) {
		points += 10
	}
	r.Points = points
	return points
}

func isAlphanumeric(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9')
}

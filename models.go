package main

import (
	"fmt"
	"strings"
	"time"
)

const (
	dateFormat = "2006-01-02"
	timeFormat = "15:04"
)

type DateTime struct {
	time.Time
}

func (dt *DateTime) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), `"`)
	if parsedDate, err := time.Parse(dateFormat, str); err == nil {
		dt.Time = parsedDate
		return nil
	}
	if parsedTime, err := time.Parse(timeFormat, str); err == nil {
		dt.Time = parsedTime
		return nil
	}
	return fmt.Errorf("invalid DateTime format: %s", str)
}

type Receipt struct {
	Retailer     string   `json:"retailer"`
	PurchaseDate DateTime `json:"purchaseDate"`
	PurchaseTime DateTime `json:"purchaseTime"`
	Items        []Item   `json:"items"`
	Total        float64  `json:"total,string"`
	Points       int      `json:"points"`
}

type Item struct {
	ShortDescription string  `json:"shortDescription"`
	Price            float64 `json:"price,string"`
}
type PointsResponse struct {
	Points int `json:"points"`
}

type ReceiptIDResponse struct {
	ID string `json:"id"`
}

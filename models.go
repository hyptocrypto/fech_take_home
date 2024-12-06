package main

import (
	"time"
)

type Receipt struct {
	Retailer     string    `json:"retailer"`
	PurchaseDate time.Time `json:"purchaseDate"`
	PurchaseTime time.Time `json:"purchaseTime"`
	Items        []Item    `json:"items"`
	Total        float64   `json:"total"`
	Points       int       `json:"points"`
}

// Used as a temp for parsing values
type StringReceipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []Item `json:"items"`
	Total        string `json:"total"`
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

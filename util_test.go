package main

import (
	"strings"
	"testing"
	"time"
)

func TestGenerateHash(t *testing.T) {
	date, _ := time.Parse("2006-01-02", "2023-12-06")
	timeOfPurchase, _ := time.Parse("15:04", "15:30")

	receipt1 := &Receipt{
		Retailer:     "RetailerA",
		PurchaseDate: date,
		PurchaseTime: timeOfPurchase,
		Items: []Item{
			{ShortDescription: "Item1", Price: 10.50},
			{ShortDescription: "Item2", Price: 5.00},
		},
		Total: 15.50,
	}

	receipt2 := &Receipt{
		Retailer:     "RetailerA",
		PurchaseDate: date,
		PurchaseTime: timeOfPurchase,
		Items: []Item{
			{ShortDescription: "Item1", Price: 10.50},
			{ShortDescription: "Item2", Price: 5.00},
		},
		Total: 15.50,
	}

	receipt3 := &Receipt{
		Retailer:     "RetailerB", // Different retailer name
		PurchaseDate: date,
		PurchaseTime: timeOfPurchase,
		Items: []Item{
			{ShortDescription: "Item1", Price: 10.50},
			{ShortDescription: "Item2", Price: 5.00},
		},
		Total: 15.50,
	}

	uuid1 := GenerateUUIDForReceipt(receipt1)
	uuid2 := GenerateUUIDForReceipt(receipt2)
	uuid3 := GenerateUUIDForReceipt(receipt3)

	// Assert that identical receipts produce the same UUID
	if uuid1 != uuid2 {
		t.Error("UUID for identical receipts should match")
	}

	// Assert that different receipts produce different UUIDs
	if uuid1 == uuid3 {
		t.Error("UUID for different receipts should not match")
	}
}

func TestCalculatePoints(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		expected int
	}{
		{
			name: "Test 1 - Round Total",
			jsonData: `{
                "retailer": "M&M Corner Market",
                "purchaseDate": "2022-03-20",
                "purchaseTime": "14:33",
                "items": [
                    {
                    "shortDescription": "Gatorade",
                    "price": "2.25"
                    },{
                    "shortDescription": "Gatorade",
                    "price": "2.25"
                    },{
                    "shortDescription": "Gatorade",
                    "price": "2.25"
                    },{
                    "shortDescription": "Gatorade",
                    "price": "2.25"
                    }
                ],
                "total": "9.00"
                }`,
			expected: 109,
		},
		{
			name: "Test 2 - 5 Items",
			jsonData: `{
                "retailer": "Target",
                "purchaseDate": "2022-01-01",
                "purchaseTime": "13:01",
                "items": [
                    {
                    "shortDescription": "Mountain Dew 12PK",
                    "price": "6.49"
                    },{
                    "shortDescription": "Emils Cheese Pizza",
                    "price": "12.25"
                    },{
                    "shortDescription": "Knorr Creamy Chicken",
                    "price": "1.26"
                    },{
                    "shortDescription": "Doritos Nacho Cheese",
                    "price": "3.35"
                    },{
                    "shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
                    "price": "12.00"
                    }
                ],
                "total": "35.35"
                }`,
			expected: 28,
		},
		{
			name: "Test 3 - 1 Item Round Total",
			jsonData: `{
                "retailer": "Target",
                "purchaseDate": "2022-01-01",
                "purchaseTime": "13:01",
                "items": [
                    {
                    "shortDescription": "Mountain Dew 12PK",
                    "price": "6.49"
                    }
                ],
                "total": "35.00"
                }`,
			expected: 87,
		},
		{
			name: "Test 4 - Round total",
			jsonData: `{
                "retailer": "M&M Corner Market",
                "purchaseDate": "2022-03-20",
                "purchaseTime": "14:33",
                "items": [
                    {
                    "shortDescription": "Gatorade1",
                    "price": "22.25"
                    },{
                    "shortDescription": "Gatorade1",
                    "price": "22.25"
                    },{
                    "shortDescription": "Gatorade1",
                    "price": "22.25"
                    },{
                    "shortDescription": "Gatorade1",
                    "price": "22.25"
                    },{
                    "shortDescription": "Gatorade1",
                    "price": "22.25"
                    },{
                    "shortDescription": "Gatorade1",
                    "price": "22.25"
                    },{
                    "shortDescription": "Gatorade1",
                    "price": "22.25"
                    },{
                    "shortDescription": "Gatorade1",
                    "price": "22.25"
                    }
                ],
                "total": "9.00"
                }`,
			expected: 159,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			receipt, err := ParseReceiptJson(strings.NewReader(test.jsonData))
			if err != nil {
				t.Error("Unmarshalling json failed")
			}
			points := CalculateReceiptPoints(receipt)

			// Verify the result matches the expected value
			if points != test.expected {
				t.Errorf("Test %q failed: expected %d, got %d", test.name, test.expected, points)
			}
		})
	}
}

curl -X POST localhost:8080/receipts/process \
    -H "Content-Type: application/json" \
    -d '{
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
                }'

package entity

import (
	"github.com/google/uuid"
	"time"
)

type Transaction struct {
	ID           string    `json:"id"`
	SellingOrder *Order    `json:"selling_order"`
	BuyingOrder  *Order    `json:"buying_order"`
	Shares       int       `json:"shares"`
	Price        float64   `json:"price"`
	Total        float64   `json:"total"`
	DateTime     time.Time `json:"date_time"`
}

func NewTransaction(sellingOrder *Order, buyingOrder *Order, shares int, price float64) *Transaction {
	total := float64(shares) * price
	return &Transaction{
		ID:           uuid.New().String(),
		SellingOrder: sellingOrder,
		BuyingOrder:  buyingOrder,
		Shares:       shares,
		Price:        price,
		Total:        total,
		DateTime:     time.Now(),
	}
}

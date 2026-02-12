package domain

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID          string
	Amount      int
	Description string
	CreatedAt   string
	UpdatedAt   string
}

func NewTransaction(amount int, description string) *Transaction {
	return &Transaction{
		ID:          uuid.NewString(),
		Amount:      amount,
		Description: description,
		CreatedAt:   time.Now().Format(time.DateTime),
		UpdatedAt:   time.Now().Format(time.DateTime),
	}
}

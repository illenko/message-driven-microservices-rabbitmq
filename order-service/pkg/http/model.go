package http

import (
	"github.com/google/uuid"
)

type OrderRequest struct {
	CustomerID uuid.UUID `json:"customerId"`
	CardID     uuid.UUID `json:"cardId"`
	ItemID     uuid.UUID `json:"itemId"`
	Price      int       `json:"price"`
}

type OrderResponse struct {
	ID         uuid.UUID `json:"id"`
	CustomerID uuid.UUID `json:"customerId"`
	CardID     uuid.UUID `json:"cardId"`
	ItemID     uuid.UUID `json:"itemId"`
	Price      int       `json:"price"`
	Status     string    `json:"status"`
	Message    string    `json:"message"`
}

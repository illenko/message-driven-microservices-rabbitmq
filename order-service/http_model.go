package main

import "github.com/google/uuid"

type OrderStatus string

const (
	Pending   OrderStatus = "pending"
	Completed OrderStatus = "completed"
	Cancelled OrderStatus = "cancelled"
)

type OrderRequest struct {
	CustomerID uuid.UUID `json:"customerId"`
	CardID     uuid.UUID `json:"cardId"`
	ItemID     uuid.UUID `json:"itemId"`
	Price      int       `json:"price"`
}

type OrderResponse struct {
	ID         uuid.UUID   `json:"id"`
	CustomerID uuid.UUID   `json:"customerId"`
	CardID     uuid.UUID   `json:"cardId"`
	ItemID     uuid.UUID   `json:"itemId"`
	Price      int         `json:"price"`
	Status     OrderStatus `json:"status"`
}

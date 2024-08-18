package amqpmodel

import "github.com/google/uuid"

type OrderAction struct {
	ID         uuid.UUID `json:"id"`
	CustomerID uuid.UUID `json:"customerId"`
	CardID     uuid.UUID `json:"cardId"`
	ItemID     uuid.UUID `json:"itemId"`
	Price      int       `json:"price"`
}

type OrderActionResultStatus string

const (
	OrderActionResultStatusSuccess OrderActionResultStatus = "success"
	OrderActionResultStatusFailed  OrderActionResultStatus = "failed"
)

type OrderActionResult struct {
	ID         uuid.UUID               `json:"id"`
	CustomerID uuid.UUID               `json:"customerId"`
	CardID     uuid.UUID               `json:"cardId"`
	ItemID     uuid.UUID               `json:"itemId"`
	Price      int                     `json:"price"`
	Status     OrderActionResultStatus `json:"status"`
}

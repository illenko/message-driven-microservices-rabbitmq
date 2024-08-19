package model

import (
	"github.com/google/uuid"
	"github.com/illenko/order-service/internal/constants"
)

type Order struct {
	ID         uuid.UUID
	CustomerID uuid.UUID
	CardID     uuid.UUID
	ItemID     uuid.UUID
	Price      int
	Status     constants.OrderStatus
	Message    constants.OrderMessage
}

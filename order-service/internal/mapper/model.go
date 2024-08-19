package mapper

import (
	"github.com/illenko/order-service/internal/constants"
	"github.com/illenko/order-service/internal/model"
	"github.com/illenko/order-service/pkg/http"
)

type OrderMapper struct{}

func NewOrderMapper() *OrderMapper {
	return &OrderMapper{}
}

func (o *OrderMapper) ToOrder(req http.OrderRequest) model.Order {
	return model.Order{
		CustomerID: req.CustomerID,
		CardID:     req.CardID,
		ItemID:     req.ItemID,
		Price:      req.Price,
		Status:     constants.StatusPending,
		Message:    constants.MessageOrderPending,
	}
}

func (o *OrderMapper) ToOrderResponse(order model.Order) http.OrderResponse {
	return http.OrderResponse{
		ID:         order.ID,
		CustomerID: order.CustomerID,
		CardID:     order.CardID,
		ItemID:     order.ItemID,
		Price:      order.Price,
		Status:     string(order.Status),
		Message:    string(order.Message),
	}
}

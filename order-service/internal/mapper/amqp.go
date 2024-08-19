package mapper

import (
	"github.com/illenko/common/amqpmodel"
	"github.com/illenko/order-service/internal/model"
)

type AmqpMapper struct{}

func NewAmqpMapper() *AmqpMapper {
	return &AmqpMapper{}
}

func (a *AmqpMapper) ToOrderAction(orderResult amqpmodel.OrderActionResult) amqpmodel.OrderAction {
	return amqpmodel.OrderAction{
		ID:         orderResult.ID,
		CustomerID: orderResult.CustomerID,
		CardID:     orderResult.CardID,
		Price:      orderResult.Price,
		ItemID:     orderResult.ItemID,
	}
}

func (a *AmqpMapper) OrderToAction(orderAction model.Order) amqpmodel.OrderAction {
	return amqpmodel.OrderAction{
		ID:         orderAction.ID,
		CustomerID: orderAction.CustomerID,
		CardID:     orderAction.CardID,
		Price:      orderAction.Price,
		ItemID:     orderAction.ItemID,
	}
}

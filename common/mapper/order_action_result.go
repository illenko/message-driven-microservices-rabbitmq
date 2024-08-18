package mapper

import "github.com/illenko/common/amqpmodel"

func ToOrderActionResult(orderAction amqpmodel.OrderAction, status amqpmodel.OrderActionResultStatus) amqpmodel.OrderActionResult {
	return amqpmodel.OrderActionResult{
		ID:         orderAction.ID,
		CustomerID: orderAction.CustomerID,
		CardID:     orderAction.CardID,
		ItemID:     orderAction.ItemID,
		Price:      orderAction.Price,
		Status:     status,
	}
}

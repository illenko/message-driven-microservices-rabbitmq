package constants

type OrderStatus string

const (
	StatusPending   OrderStatus = "pending"
	StatusCompleted OrderStatus = "completed"
	StatusFailed    OrderStatus = "failed"
)

type OrderMessage string

const (
	MessageOrderPending       OrderMessage = "Order is pending"
	MessageOrderCompleted     OrderMessage = "Order is completed"
	MessageOrderFailed        OrderMessage = "Order is failed"
	MessageOrderFailedItem    OrderMessage = "Order is failed, item is out of stock"
	MessageOrderFailedPayment OrderMessage = "Order is failed, payment failed"
)

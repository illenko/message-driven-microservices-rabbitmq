package messaging

type QueueConfig struct {
	Name string
	TTL  int
	DLX  string
}

type Config struct {
	ConnectionUrl                 string
	OrderActionExchange           string
	OrderResultExchange           string
	DLXExchange                   string
	ProductReservationQueue       QueueConfig
	ProductCancellationQueue      QueueConfig
	PaymentQueue                  QueueConfig
	ProductReservationResultQueue QueueConfig
	PaymentResultQueue            QueueConfig
	DLXProductReservationQueue    QueueConfig
	DLXPaymentQueue               QueueConfig
}

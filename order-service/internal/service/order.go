package service

import (
	"log"

	"github.com/google/uuid"
	"github.com/illenko/common/amqpmodel"
	"github.com/illenko/order-service/internal/constants"
	"github.com/illenko/order-service/internal/messaging"
	"github.com/illenko/order-service/internal/model"
	"github.com/illenko/order-service/pkg/http"
)

type AmqpMapper interface {
	ToOrderAction(orderResult amqpmodel.OrderActionResult) amqpmodel.OrderAction
	OrderToAction(orderAction model.Order) amqpmodel.OrderAction
}

type OrderMapper interface {
	ToOrder(req http.OrderRequest) model.Order
	ToOrderResponse(order model.Order) http.OrderResponse
}

type OrderRepository interface {
	Create(order model.Order) (model.Order, error)
	Update(order model.Order) (model.Order, error)
	FindById(id uuid.UUID) (model.Order, bool)
}

type OrderActionPublisher interface {
	PublishOrderAction(exchangeName, routingKey string, order amqpmodel.OrderAction) error
}

type OrderService struct {
	messagingConfig messaging.Config
	publisher       OrderActionPublisher
	repo            OrderRepository
	orderMapper     OrderMapper
	amqpMapper      AmqpMapper
}

func NewOrderService(messagingConfig messaging.Config, publisher OrderActionPublisher, repo OrderRepository, orderMapper OrderMapper, amqpMapper AmqpMapper) *OrderService {
	return &OrderService{
		messagingConfig: messagingConfig,
		publisher:       publisher,
		repo:            repo,
		orderMapper:     orderMapper,
		amqpMapper:      amqpMapper,
	}
}

func (s *OrderService) CreateOrder(request http.OrderRequest) (http.OrderResponse, error) {
	order, err := s.repo.Create(s.orderMapper.ToOrder(request))

	if err != nil {
		log.Printf("Failed to save order: %v", err)
		return http.OrderResponse{}, err
	}

	orderResponse := s.orderMapper.ToOrderResponse(order)
	log.Printf("Order created: %v", orderResponse)

	err = s.publishOrderAction(s.messagingConfig.OrderActionExchange, s.messagingConfig.ProductReservationQueue.Name, s.amqpMapper.OrderToAction(order))

	if err != nil {
		log.Printf("Failed to publish product reservation action: %v", err)
		return http.OrderResponse{}, err
	}

	return orderResponse, nil
}

func (s *OrderService) GetOrder(id uuid.UUID) (http.OrderResponse, bool) {
	log.Printf("Getting order with ID: %v", id)
	order, exists := s.repo.FindById(id)
	if !exists {
		log.Printf("Order with ID %v not found", id)
		return http.OrderResponse{}, false
	}

	orderResponse := s.orderMapper.ToOrderResponse(order)
	log.Printf("Order found: %v", orderResponse)
	return orderResponse, true
}

func (s *OrderService) ProcessExpiredProductReservation(orderAction amqpmodel.OrderAction) {
	log.Printf("Processing expired product reservation: %v", orderAction)
	s.updateOrder(orderAction.ID, constants.StatusFailed, constants.MessageOrderFailed)
}

func (s *OrderService) ProcessExpiredPayment(orderAction amqpmodel.OrderAction) {
	log.Printf("Processing expired payment: %v", orderAction)
	s.updateOrder(orderAction.ID, constants.StatusFailed, constants.MessageOrderFailed)

	err := s.publishOrderAction(s.messagingConfig.OrderActionExchange, s.messagingConfig.ProductCancellationQueue.Name, orderAction)
	if err != nil {
		log.Printf("Failed to publish product cancellation action: %v", err)
		s.updateOrder(orderAction.ID, constants.StatusFailed, constants.MessageOrderFailed)
		return
	}

	log.Printf("Published product cancellation action: %v", orderAction)
}

func (s *OrderService) ProcessProductReservationResult(orderResult amqpmodel.OrderActionResult) {
	log.Printf("Processing product reservation result: %v", orderResult)

	if orderResult.Status == amqpmodel.OrderActionResultStatusFailed {
		log.Printf("Product reservation failed: %v", orderResult)
		s.updateOrder(orderResult.ID, constants.StatusFailed, constants.MessageOrderFailedItem)
		return
	}

	err := s.publishOrderAction(s.messagingConfig.OrderActionExchange, s.messagingConfig.PaymentQueue.Name, s.amqpMapper.ToOrderAction(orderResult))
	if err != nil {
		log.Printf("Failed to publish payment action: %v", err)
		s.updateOrder(orderResult.ID, constants.StatusFailed, constants.MessageOrderFailed)
		return
	}

	log.Printf("Published payment action: %v", orderResult)
}

func (s *OrderService) ProcessPaymentResult(orderResult amqpmodel.OrderActionResult) {
	log.Printf("Processing payment result: %v", orderResult)

	if orderResult.Status == amqpmodel.OrderActionResultStatusFailed {
		log.Printf("Payment failed: %v", orderResult)
		s.updateOrder(orderResult.ID, constants.StatusFailed, constants.MessageOrderFailedPayment)

		err := s.publishOrderAction(s.messagingConfig.OrderActionExchange, s.messagingConfig.ProductCancellationQueue.Name, s.amqpMapper.ToOrderAction(orderResult))
		if err != nil {
			log.Printf("Failed to publish product cancellation action: %v", err)
			s.updateOrder(orderResult.ID, constants.StatusFailed, constants.MessageOrderFailed)
			return
		}

		log.Printf("Published product cancellation action: %v", orderResult)
		return
	}

	s.updateOrder(orderResult.ID, constants.StatusCompleted, constants.MessageOrderCompleted)
	log.Printf("OrderService completed: %v", orderResult)
}

func (s *OrderService) publishOrderAction(exchangeName, routingKey string, order amqpmodel.OrderAction) error {
	log.Printf("Publishing order action: %v to exchange: %v with routing key: %v", order, exchangeName, routingKey)
	err := s.publisher.PublishOrderAction(exchangeName, routingKey, order)
	if err != nil {
		log.Printf("Failed to publish order action: %v", err)
		return err
	}
	return nil
}

func (s *OrderService) updateOrder(id uuid.UUID, status constants.OrderStatus, message constants.OrderMessage) {
	log.Printf("Updating order with ID: %v", id)
	order, exists := s.repo.FindById(id)
	if !exists {
		log.Printf("Order with ID %v not found", id)
		return
	}

	order.Status = status
	order.Message = message

	order, err := s.repo.Update(order)
	if err != nil {
		log.Printf("Failed to update order: %v", err)
	}
	log.Printf("Order updated: %v", order)
}

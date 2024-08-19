package repository

import (
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/illenko/order-service/internal/model"
)

type InMemoryOrderRepository struct {
	orders map[uuid.UUID]model.Order
	mu     sync.Mutex
}

func NewInMemoryOrderRepository() *InMemoryOrderRepository {
	return &InMemoryOrderRepository{
		orders: make(map[uuid.UUID]model.Order),
	}
}

func (r *InMemoryOrderRepository) Create(order model.Order) (model.Order, error) {
	order.ID = uuid.New()
	log.Printf("Creating order with ID: %v", order.ID)
	return r.Update(order)
}

func (r *InMemoryOrderRepository) Update(order model.Order) (model.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.orders[order.ID] = order
	log.Printf("Updated order with ID: %v", order.ID)
	return order, nil
}

func (r *InMemoryOrderRepository) FindById(id uuid.UUID) (model.Order, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	order, exists := r.orders[id]
	if exists {
		log.Printf("Found order with ID: %v", id)
	} else {
		log.Printf("Order with ID: %v not found", id)
	}
	return order, exists
}

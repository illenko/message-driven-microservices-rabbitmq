package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/illenko/order-service/internal/service"
	httpModel "github.com/illenko/order-service/pkg/http"
)

type OrderService interface {
	CreateOrder(request httpModel.OrderRequest) (httpModel.OrderResponse, error)
	GetOrder(id uuid.UUID) (httpModel.OrderResponse, bool)
}

type OrderHandler interface {
	CreateOrder(w http.ResponseWriter, r *http.Request)
	GetOrder(w http.ResponseWriter, r *http.Request)
}

type OrderHandlerImpl struct {
	orderService *service.OrderService
}

func NewOrderHandlerImpl(orderService *service.OrderService) *OrderHandlerImpl {
	return &OrderHandlerImpl{orderService: orderService}
}

func (o *OrderHandlerImpl) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var orderReq httpModel.OrderRequest
	if err := json.NewDecoder(r.Body).Decode(&orderReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	order, err := o.orderService.CreateOrder(orderReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(order)
}

func (o *OrderHandlerImpl) GetOrder(w http.ResponseWriter, r *http.Request) {
	orderID := r.PathValue("id")
	if orderID == "" {
		http.Error(w, "order_id is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(orderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	order, exists := o.orderService.GetOrder(id)
	if !exists {
		http.Error(w, "order not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

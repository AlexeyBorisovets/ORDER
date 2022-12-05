package service

import (
	"github.com/AlexeyBorisovets/ORDER/internal/model"
	"github.com/AlexeyBorisovets/ORDER/internal/repository"

	"context"
)

// Service s
type Service struct {
	repos repository.Repository
}

// NewService create new service connection
func NewService(pool repository.Repository) *Service {
	return &Service{repos: pool}
}

func (se *Service) CreateOrder(ctx context.Context, order *model.Order) (string, error) {
	return se.repos.CreateOrder(ctx, order)
}

func (se *Service) GetOrder(ctx context.Context, orderID string) (*model.Order, error) {
	return se.repos.GetOrderByID(ctx, orderID)
}

func (se *Service) GetOrdersByConsID(ctx context.Context, consID string) ([]*model.Order, error) {
	return se.repos.GetOrdersByConsID(ctx, consID)
}

func (se *Service) GetOrdersByVendID(ctx context.Context, vendID string) ([]*model.Order, error) {
	return se.repos.GetOrdersByVendID(ctx, vendID)
}

func (se *Service) UpdateOrder(ctx context.Context, orderid string, order *model.Order) error {
	return se.repos.UpdateOrder(ctx, orderid, order)
}

func (se *Service) DeleteUser(ctx context.Context, orderid string) error {
	return se.DeleteUser(ctx, orderid)
}

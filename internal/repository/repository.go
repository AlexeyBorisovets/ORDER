package repository

import (
	"context"

	"github.com/AlexeyBorisovets/ORDER/internal/model"
)

type Repository interface {
	CreateOrder(ctx context.Context, order *model.Order) (string, error)
	GetOrderByID(ctx context.Context, orderID string) (*model.Order, error)
	GetOrdersByConsID(ctx context.Context, consID string) ([]*model.Order, error)
	GetOrdersByVendID(ctx context.Context, vendID string) ([]*model.Order, error)
	UpdateOrder(ctx context.Context, orderid string, order *model.Order) error
	DeleteUser(ctx context.Context, orderid string) error
}

package repository

import (
	"context"
	"fmt"

	"github.com/AlexeyBorisovets/ORDER/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/gommon/log"
)

type PRepository struct {
	Pool *pgxpool.Pool
}

func (p *PRepository) CreateOrder(ctx context.Context, order *model.Order) (string, error) {
	orderID := uuid.New().String()
	_, err := p.Pool.Exec(ctx, "insert into orders(orderid,productid,consumerid,vendorid,amount,orderprice) values($1,$2,$3,$4,$5,$6)",
		orderID, &order.ProductID, &order.ConsumerID, &order.VendorID, &order.Amount, &order.OrderPrice)
	if err != nil {
		log.Errorf("database error with create order: %v", err)
		return "", err
	}
	return orderID, nil
}

func (p *PRepository) GetOrderByID(ctx context.Context, orderID string) (*model.Order, error) {
	order := model.Order{}
	err := p.Pool.QueryRow(ctx, "select orderid,productid,consumerid,vendorid,amount,orderprice from orders where orderid=$1", orderID).Scan(
		&order.OrderID, &order.ProductID, &order.ConsumerID, &order.VendorID, &order.Amount, &order.OrderPrice)
	if err != nil {
		if err == pgx.ErrNoRows {
			return &model.Order{}, fmt.Errorf("error: no such order in db: %v", err)
		}
		log.Errorf("database error,GetOrderByID: %v", err)
		return &model.Order{}, err
	}
	return &order, nil
}

func (p *PRepository) GetOrdersByConsID(ctx context.Context, consID string) ([]*model.Order, error) {
	var orders []*model.Order
	orders_db, err := p.Pool.Query(ctx, "select orderid,productid,consumerid,vendorid,amount,orderprice from orders where consumerid=$1", consID)
	if err != nil {
		log.Errorf("database error with find orders from consID, %v", err)
		return nil, err
	}
	defer orders_db.Close()
	for orders_db.Next() {
		order := model.Order{}
		err = orders_db.Scan(&order.OrderID, &order.ProductID, &order.ConsumerID, &order.VendorID, &order.Amount, &order.OrderPrice)
		if err != nil {
			log.Errorf("database error with find orders from consID, %v", err)
			return nil, err
		}
		orders = append(orders, &order)
	}

	return orders, nil
}

func (p *PRepository) GetOrdersByVendID(ctx context.Context, vendID string) ([]*model.Order, error) {
	var orders []*model.Order
	orders_db, err := p.Pool.Query(ctx, "select orderid,productid,consumerid,vendorid,amount,orderprice from orders where consumerid=$1", vendID)
	if err != nil {
		log.Errorf("database error with find orders from vendID, %v", err)
		return nil, err
	}
	defer orders_db.Close()
	for orders_db.Next() {
		order := model.Order{}
		err = orders_db.Scan(&order.OrderID, &order.ProductID, &order.ConsumerID, &order.VendorID, &order.Amount, &order.OrderPrice)
		if err != nil {
			log.Errorf("database error with find orders from vendID, %v", err)
			return nil, err
		}
		orders = append(orders, &order)
	}

	return orders, nil
}

func (p *PRepository) UpdateOrder(ctx context.Context, orderid string, order *model.Order) error {
	a, err := p.Pool.Exec(ctx, "update orders set orderid=$1,productid=$2,consumerid=$3,vendorid=$4,amount=$5,orderprice=$6 where orderid=$7",
		&order.OrderID, &order.ProductID, &order.ConsumerID, &order.VendorID, &order.Amount, &order.OrderPrice, orderid)
	if a.RowsAffected() == 0 {
		return fmt.Errorf("no order with such id")
	}
	if err != nil {
		log.Errorf("error with update order %v", err)
		return err
	}
	return nil
}

func (p *PRepository) DeleteUser(ctx context.Context, orderid string) error {
	a, err := p.Pool.Exec(ctx, "delete from orders where id=$1", orderid)
	if a.RowsAffected() == 0 {
		return fmt.Errorf("no order with such id")
	}
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("no order with such id: %v", err)
		}
		log.Errorf("error with delete order %v", err)
		return err
	}
	return nil
}

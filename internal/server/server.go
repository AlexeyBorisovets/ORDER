package server

import (
	"fmt"

	"github.com/AlexeyBorisovets/ORDER/internal/model"
	"github.com/AlexeyBorisovets/ORDER/internal/service"

	proto_order "github.com/AlexeyBorisovets/ORDER/proto"
	proto_user "github.com/AlexeyBorisovets/USER/proto"
	proto_basket "github.com/VladZheleznov/shopping-basket/proto"

	"context"
)

type Server struct {
	uc proto_user.USERClient
	bc proto_basket.CRUDClient
	se *service.Service
	proto_order.ORDERServer
}

func NewServer(serv *service.Service, new_uc proto_user.USERClient, new_bc proto_basket.CRUDClient) *Server {
	return &Server{se: serv, uc: new_uc, bc: new_bc}
}

func (s *Server) CreateOrder(ctx context.Context, request *proto_order.CreateOrderRequest) (*proto_order.CreateOrderResponse, error) {
	resp_c_balance, err := s.uc.GetBalanceByID(ctx, &proto_user.GetBalanceByIDRequest{Id: request.Consumerid})
	if err != nil {
		return nil, fmt.Errorf("server: error while connect to USER_SERVICE, %e", err)
	}
	resp_v_balance, err := s.uc.GetBalanceByID(ctx, &proto_user.GetBalanceByIDRequest{Id: request.Vendorid})
	if err != nil {
		return nil, fmt.Errorf("server: error while connect to USER_SERVICE, %e", err)
	}
	resp_product, err := s.bc.GetItem(ctx, &proto_basket.GetItemRequest{Id: request.Productid})
	if err != nil {
		return nil, fmt.Errorf("server: error while connect to basket_SERVICE, %e", err)
	}
	cons_balance := resp_c_balance.GetBalance()
	vendor_balance := resp_v_balance.GetBalance()
	order_price := request.Orderprice
	stock_amount := resp_product.Product.Quantity
	order_amount := request.Amount

	if cons_balance >= order_price && stock_amount >= int32(order_amount) {
		o := model.Order{
			ProductID:  request.Productid,
			ConsumerID: request.Consumerid,
			VendorID:   request.Vendorid,
			Amount:     uint(request.Amount),
			OrderPrice: uint(request.Orderprice),
		}
		new_order_id, err := s.se.CreateOrder(ctx, &o)
		if err != nil {
			return nil, err
		}
		new_cons_balance := cons_balance - order_price
		new_vend_balance := vendor_balance + order_price
		_, err = s.uc.UpdateBalance(ctx, &proto_user.UpdateBalanceRequest{Id: request.Consumerid, Balance: new_cons_balance})
		if err != nil {
			return nil, fmt.Errorf("server: error while connect to USER_SERVICE, %e", err)
		}
		_, err = s.uc.UpdateBalance(ctx, &proto_user.UpdateBalanceRequest{Id: request.Vendorid, Balance: new_vend_balance})
		if err != nil {
			return nil, fmt.Errorf("server: error while connect to USER_SERVICE, %e", err)
		}
		new_stock_amount := stock_amount - int32(order_amount)
		p := proto_basket.Product{Quantity: new_stock_amount}
		_, err = s.bc.UpdateItem(ctx, &proto_basket.UpdateItemRequest{Id: request.Productid, Product: &p})
		if err != nil {
			return nil, fmt.Errorf("server: error while connect to USER_SERVICE, %e", err)
		}
		return &proto_order.CreateOrderResponse{Orderid: new_order_id}, nil
	} else {
		return nil, fmt.Errorf("server: not enough money or not enough products at stock, %e", err)
	}
}

func (s *Server) UpdateOrder(ctx context.Context, request *proto_order.UpdateOrderRequest) (*proto_order.Response, error) {
	order := &model.Order{
		ProductID:  request.Order.Productid,
		VendorID:   request.Order.Vendorid,
		ConsumerID: request.Order.Consumerid,
		Amount:     uint(request.Order.Amount),
		OrderPrice: uint(request.Order.Orderprice),
	}
	OrderID := request.GetOrderid()
	err := s.se.UpdateOrder(ctx, OrderID, order)
	if err != nil {
		return nil, fmt.Errorf("server: error during updating order, %e", err)
	}
	return new(proto_order.Response), nil
}

func (s *Server) DeleteOrder(ctx context.Context, request *proto_order.DeleteOrderRequest) (*proto_order.Response, error) {
	Order_id := request.GetOrderid()
	err := s.se.DeleteUser(ctx, Order_id)
	if err != nil {
		return nil, fmt.Errorf("server: error during deletion order, %e", err)
	}
	return new(proto_order.Response), nil
}

func (s *Server) GetOrder(ctx context.Context, request *proto_order.GetOrderRequest) (*proto_order.GetOrderResponse, error) {
	orderID := request.GetOrderid()
	order, err := s.se.GetOrder(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("server: error during getting order from DB, %e", err)
	}
	newOrder := &proto_order.Order{
		Orderid:    request.Orderid,
		Productid:  order.ProductID,
		Vendorid:   order.VendorID,
		Consumerid: order.ConsumerID,
		Amount:     uint64(order.Amount),
		Orderprice: uint64(order.OrderPrice),
	}
	return &proto_order.GetOrderResponse{Order: newOrder}, nil
}

func (s *Server) GetOrderByConsIDRequest(ctx context.Context, request *proto_order.GetOrderByConsIDRequest) (*proto_order.GetOrderByConsIDResponse, error) {
	consID := request.GetConsumerid()
	orders, err := s.se.GetOrdersByConsID(ctx, consID)
	if err != nil {
		return nil, fmt.Errorf("server: error during getting ordersbyconsID from DB, %e", err)
	}
	var buf []*proto_order.Order
	for _, order := range orders {
		order_proto := new(proto_order.Order)
		order_proto.Orderid = order.OrderID
		order_proto.Productid = order.ProductID
		order_proto.Vendorid = order.VendorID
		order_proto.Consumerid = order.ConsumerID
		order_proto.Amount = uint64(order.Amount)
		order_proto.Orderprice = uint64(order.OrderPrice)
		buf = append(buf, order_proto)
	}
	return &proto_order.GetOrderByConsIDResponse{Order: buf}, nil
}

func (s *Server) GetOrderByVendIDRequest(ctx context.Context, request *proto_order.GetOrderByVendIDRequest) (*proto_order.GetOrderByVendIDResponse, error) {
	consID := request.GetVendorid()
	orders, err := s.se.GetOrdersByVendID(ctx, consID)
	if err != nil {
		return nil, fmt.Errorf("server: error during getting ordersbyvendID from DB, %e", err)
	}
	var buf []*proto_order.Order
	for _, order := range orders {
		order_proto := new(proto_order.Order)
		order_proto.Orderid = order.OrderID
		order_proto.Productid = order.ProductID
		order_proto.Vendorid = order.VendorID
		order_proto.Consumerid = order.ConsumerID
		order_proto.Amount = uint64(order.Amount)
		order_proto.Orderprice = uint64(order.OrderPrice)
		buf = append(buf, order_proto)
	}
	return &proto_order.GetOrderByVendIDResponse{Order: buf}, nil
}

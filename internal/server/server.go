package server

import (
	"fmt"

	"github.com/AlexeyBorisovets/ORDER/internal/model"
	"github.com/AlexeyBorisovets/ORDER/internal/service"

	proto_order "github.com/AlexeyBorisovets/ORDER/internal/proto"
	proto_user "github.com/AlexeyBorisovets/USER/proto"
	proto_basket "github.com/VladZheleznov/shopping-basket/proto"
	model_basket "github.com/VladZheleznov/shopping-basket/internal/model"

	"context"
)

type Server struct {
	uc proto_user.USERClient
	bc proto_basket.CRUDClient
	se *service.Service
	proto_order.ORDERServer
}

func NewServer(serv *service.Service, new_uc proto_user.USERClient,new_bc proto_basket.CRUDClient) *Server {
	return &Server{se: serv, uc: new_uc, bc: new_bc}
}

func (s *Server) CreateOrder(ctx context.Context, request *proto_order.CreateOrderRequest) (*proto_order.CreateOrderResponse, error) {
	resp_c_balance, err := s.uc.GetBalanceByID(ctx, &proto_user.GetBalanceByIDRequest{Id: request.Consumerid})
	if err != nil{
		return nil,fmt.Errorf("server: error while connect to USER_SERVICE, %e", err)
	}
	resp_v_balance, err := s.uc.GetBalanceByID(ctx, &proto_user.GetBalanceByIDRequest{Id: request.Vendorid})
	if err != nil{
		return nil,fmt.Errorf("server: error while connect to USER_SERVICE, %e", err)
	}
	resp_product, err := s.bc.GetItem(ctx,&proto_basket.GetItemRequest{Id: request.Productid})
	if err != nil{
		return nil,fmt.Errorf("server: error while connect to basket_SERVICE, %e", err)
	}
	cons_balance := resp_c_balance.GetBalance()
	vendor_balance := resp_v_balance.GetBalance()
	order_price := request.Orderprice
	stock_amount := resp_product.Product.Quantity
	order_amount := request.Amount

	if cons_balance >= order_price && stock_amount >= int32(order_amount){
		o := model.Order{
			ProductID: request.Productid,
			ConsumerID: request.Consumerid,
			VendorID: request.Vendorid,
			Amount: uint(request.Amount),
			OrderPrice: uint(request.Orderprice),
		}
		new_order_id , err := s.se.CreateOrder(ctx,&o)
		if err != nil{
			return nil,err
		}
		new_cons_balance := cons_balance - order_price
		new_vend_balance := vendor_balance + order_price
		_,err = s.uc.UpdateBalance(ctx,&proto_user.UpdateBalanceRequest{Id: request.Consumerid,Balance: new_cons_balance})
		if err != nil{
			return nil,fmt.Errorf("server: error while connect to USER_SERVICE, %e", err)
		}
		_,err = s.uc.UpdateBalance(ctx,&proto_user.UpdateBalanceRequest{Id: request.Vendorid,Balance: new_vend_balance})
		if err != nil{
			return nil,fmt.Errorf("server: error while connect to USER_SERVICE, %e", err)
		}
		new_stock_amount := stock_amount - int32(order_amount)
		p := model_basket.Product{

		}
		_,err = s.bc.UpdateItem(ctx,&proto_basket.UpdateItemRequest{Id: request.Productid,}) 


	}


	
	

}

package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/AlexeyBorisovets/ORDER/internal/repository"
	"github.com/AlexeyBorisovets/ORDER/internal/server"
	"github.com/AlexeyBorisovets/ORDER/internal/service"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"

	proto_order "github.com/AlexeyBorisovets/ORDER/proto"
	proto_user "github.com/AlexeyBorisovets/USER/proto"
	proto_basket "github.com/VladZheleznov/shopping-basket/proto"
)

func main() {
	ctx := context.Background()
	listen, err := net.Listen("tcp", ":8082")
	if err != nil {
		defer log.Fatalf("error with starting server: %e", err)
	}
	fmt.Println("Server successfully started on port :8082...")

	poolP, err := pgxpool.Connect(ctx, "postgres://postgres:123@localhost:5432/Test")
	if err != nil {
		log.Fatalf("error with starting postgresql: %v", err)
	} else {
		fmt.Println("DB successfully connect...")
	}
	rP := repository.PRepository{Pool: poolP}
	defer func() {
		poolP.Close()
	}()

	conn_us, err := grpc.Dial(":8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("error with connect to user service: %v", err)
	}
	conn_ps, err := grpc.Dial(":8081", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("error with connect to product service: %v", err)
	}

	user_client := proto_user.NewUSERClient(conn_us)
	product_client := proto_basket.NewCRUDClient(conn_ps)

	grpc_ns := grpc.NewServer()
	new_service := service.NewService(&rP)
	order_ns := server.NewServer(new_service, user_client, product_client)
	proto_order.RegisterORDERServer(grpc_ns, order_ns)

	if err = grpc_ns.Serve(listen); err != nil {
		defer log.Fatalf("error while listening server: %e", err)
	}

}

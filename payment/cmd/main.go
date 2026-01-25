package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	paymentv1 "github.com/qyrlabs/test-backend/shared/pkg/proto/payment/v1"
)

const grpcPort = 50062

type paymentService struct {
	paymentv1.UnimplementedPaymentServiceServer
}

// Initiates order payment.
func (s *paymentService) PayOrder(ctx context.Context, req *paymentv1.PayOrderRequest) (*paymentv1.PayOrderResponse, error) {
	if _, err := uuid.Parse(req.GetOrderUuid()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid order_uuid format: %v", err)
	}
	if _, err := uuid.Parse(req.GetUserUuid()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user_uuid format: %v", err)
	}
	if req.GetPaymentMethod() == 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid payment method")
	}

	uuid := uuid.New().String()
	log.Printf("Payment succeed, transaction_uuid: %s", uuid)
	return &paymentv1.PayOrderResponse{
		TransactionUuid: uuid,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}

	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v\n", cerr)
		}
	}()

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	service := &paymentService{}

	paymentv1.RegisterPaymentServiceServer(grpcServer, service)

	go func() {
		log.Printf("gRPC server listening on %s\n", lis.Addr().String())
		err := grpcServer.Serve(lis)
		if err != nil {
			log.Printf("Failed to serve: %v\n", err)
			return
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down gRPC server...")
	grpcServer.GracefulStop()
	log.Println("gRPC server stopped")
}

package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	apiinventoryv1 "github.com/qyrlabs/test-backend/inventory/internal/api/inventory/v1"
	partRepository "github.com/qyrlabs/test-backend/inventory/internal/repository/part"
	partService "github.com/qyrlabs/test-backend/inventory/internal/service/part"
	protoinventoryv1 "github.com/qyrlabs/test-backend/shared/pkg/proto/inventory/v1"
)

const grpcPort = 50052

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	repo := partRepository.NewRepository()
	service := partService.NewService(repo)
	api := apiinventoryv1.NewAPI(service)

	protoinventoryv1.RegisterInventoryServiceServer(grpcServer, api)

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

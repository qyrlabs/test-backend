package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	inventoryv1 "github.com/qyrlabs/test-backend/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcPort = 50061

type inventoryService struct {
	inventoryv1.UnimplementedInventoryServiceServer

	mu    sync.RWMutex
	parts map[string]*inventoryv1.Part
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

	service := &inventoryService{
		parts: make(map[string]*inventoryv1.Part),
	}

	inventoryv1.RegisterInventoryServiceServer(grpcServer, service)

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

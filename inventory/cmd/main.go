package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"slices"
	"sync"
	"syscall"

	"github.com/brianvoe/gofakeit/v7"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	inventoryv1 "github.com/qyrlabs/test-backend/shared/pkg/proto/inventory/v1"
)

const grpcPort = 50061

type inventoryService struct {
	inventoryv1.UnimplementedInventoryServiceServer

	mu    sync.RWMutex
	parts map[string]*inventoryv1.Part
}

// Get part info by its UUID.
func (s *inventoryService) GetPart(ctx context.Context, req *inventoryv1.GetPartRequest) (*inventoryv1.GetPartResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	part, ok := s.parts[req.GetUuid()]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "part with uuid %s is not found", req.GetUuid())
	}

	return &inventoryv1.GetPartResponse{
		Part: part,
	}, nil
}

// Returns List of Parts by filter.
func (s *inventoryService) ListParts(ctx context.Context, req *inventoryv1.ListPartsRequest) (*inventoryv1.ListPartsResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	filter := req.GetFilter()
	uuids := filter.GetUuids()
	names := filter.GetNames()
	categories := filter.GetCategories()
	countries := filter.GetManufacturerCountries()
	tags := filter.GetTags()

	filteredParts := make([]*inventoryv1.Part, 0)
	for uuid, part := range s.parts {
		if (uuids == nil || slices.Contains(uuids, uuid)) &&
			(names == nil || slices.Contains(names, part.GetName())) &&
			(categories == nil || slices.Contains(categories, part.GetCategory())) &&
			(countries == nil || slices.Contains(countries, part.GetManufacturer().GetCountry())) &&
			(tags == nil || slices.Equal(tags, part.GetTags())) {

			filteredParts = append(filteredParts, part)
		}
	}

	if len(filteredParts) == 0 {
		return nil, status.Errorf(codes.NotFound, "parts not found")
	}

	return &inventoryv1.ListPartsResponse{
		Parts: filteredParts,
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

	service := &inventoryService{
		parts: initParts(100),
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

func initParts(count int) map[string]*inventoryv1.Part {
	parts := createParts(count)
	partsMap := make(map[string]*inventoryv1.Part, count)
	for _, part := range parts {
		partsMap[part.GetUuid()] = part
	}
	return partsMap
}

func fakeDimensions() *inventoryv1.Dimensions {
	return &inventoryv1.Dimensions{
		Length: gofakeit.Float64Range(1.0, 300.0),
		Width:  gofakeit.Float64Range(1.0, 300.0),
		Height: gofakeit.Float64Range(0.5, 150.0),
		Weight: gofakeit.Float64Range(0.1, 500.0),
	}
}

func fakeManufacturer() *inventoryv1.Manufacturer {
	return &inventoryv1.Manufacturer{
		Name:    gofakeit.Company(),
		Country: gofakeit.Country(),
		Website: gofakeit.URL(),
	}
}

func fakeTags() []string {
	tags := make([]string, 0, 5)
	for range gofakeit.IntRange(1, 5) {
		tags = append(tags, gofakeit.Word())
	}
	return tags
}

func randomCategory() inventoryv1.Category {
	// Ignore any, but not UNSPECIFIED
	vals := []inventoryv1.Category{
		inventoryv1.Category_CATEGORY_ENGINE,
		inventoryv1.Category_CATEGORY_FUEL,
		inventoryv1.Category_CATEGORY_PORTHOLE,
		inventoryv1.Category_CATEGORY_WING,
	}
	return vals[gofakeit.IntRange(0, len(vals)-1)]
}

func createParts(count int) []*inventoryv1.Part {
	parts := make([]*inventoryv1.Part, 0, count)
	for range count {
		parts = append(parts, &inventoryv1.Part{
			Uuid:          gofakeit.UUID(),
			Name:          gofakeit.Name(),
			Description:   gofakeit.Sentence(10),
			PriceMinor:    int64(gofakeit.IntRange(1, 100000)),
			StockQuantity: int64(gofakeit.IntRange(1, 100)),
			Category:      randomCategory(),
			Dimensions:    fakeDimensions(),
			Manufacturer:  fakeManufacturer(),
			Tags:          fakeTags(),
			CreatedAt:     timestamppb.New(gofakeit.Date()),
			UpdatedAt:     timestamppb.New(gofakeit.Date()),
		})
	}
	return parts
}

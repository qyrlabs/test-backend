package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	orderv1 "github.com/qyrlabs/test-backend/shared/pkg/openapi/order/v1"
	inventoryv1 "github.com/qyrlabs/test-backend/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/qyrlabs/test-backend/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	httpPort                = "8080"
	inventoryServiceAddress = "localhost:50061"
	paymentServiceAddress   = "localhost:50062"

	// Timeouts for HTTP-Server
	requestTimeout    = 10 * time.Second
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

// Repo

type OrderStorage struct {
	mutex  sync.RWMutex
	orders map[string]*orderv1.Order
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[string]*orderv1.Order),
	}
}

func (s *OrderStorage) GetOrder(uuid string) *orderv1.Order {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.orders[uuid]
}

func (s *OrderStorage) UpdateOrder(uuid string, order *orderv1.Order) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.orders[uuid] = order
}

// Handler

type OrderHandler struct {
	storage         *OrderStorage
	inventoryClient inventoryv1.InventoryServiceClient
	paymentClient   paymentv1.PaymentServiceClient
}

func NewOrderHandler(storage *OrderStorage, inventoryClient inventoryv1.InventoryServiceClient, paymentClient paymentv1.PaymentServiceClient) *OrderHandler {
	return &OrderHandler{
		storage:         storage,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
	}
}

// CancelOrder implements cancelOrder operation.
//
// Cancels an existing order.
//
// POST /api/v1/orders/{order_uuid}/cancel
func (h *OrderHandler) CancelOrder(ctx context.Context, params orderv1.CancelOrderParams) (orderv1.CancelOrderRes, error) {
	return nil, &orderv1.GenericErrorStatusCode{
		StatusCode: http.StatusNoContent,
		Response:   orderv1.GenericError{},
	}
}

// CreateOrder implements createOrder operation.
//
// Creates a new order.
//
// POST /api/v1/orders
func (h *OrderHandler) CreateOrder(ctx context.Context, req *orderv1.OrderCreateRequest) (orderv1.CreateOrderRes, error) {
	return nil, &orderv1.GenericErrorStatusCode{
		StatusCode: http.StatusNoContent,
		Response:   orderv1.GenericError{},
	}
}

// GetOrderByUuid implements getOrderByUuid operation.
//
// Retrieves order details by UUID.
//
// GET /api/v1/orders/{order_uuid}
func (h *OrderHandler) GetOrderByUuid(ctx context.Context, params orderv1.GetOrderByUuidParams) (orderv1.GetOrderByUuidRes, error) {
	return nil, &orderv1.GenericErrorStatusCode{
		StatusCode: http.StatusNoContent,
		Response:   orderv1.GenericError{},
	}
}

// PayOrder implements payOrder operation.
//
// Processes payment for an existing order.
//
// POST /api/v1/orders/{order_uuid}/pay
func (h *OrderHandler) PayOrder(ctx context.Context, req *orderv1.OrderPayRequest, params orderv1.PayOrderParams) (orderv1.PayOrderRes, error) {
	return nil, &orderv1.GenericErrorStatusCode{
		StatusCode: http.StatusNoContent,
		Response:   orderv1.GenericError{},
	}
}

// NewError creates *GenericErrorStatusCode from error returned by handler.
//
// Used for common default response.
func (h *OrderHandler) NewError(ctx context.Context, err error) *orderv1.GenericErrorStatusCode {
	return &orderv1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: orderv1.GenericError{
			Code:    orderv1.NewOptInt(http.StatusInternalServerError),
			Message: orderv1.NewOptString(err.Error()),
		},
	}
}

func initApplication() (*grpc.ClientConn, *grpc.ClientConn, *orderv1.Server, error) {
	inventoryConn, err := grpc.NewClient(
		inventoryServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create inventory service grpc connection: %w", err)
	}

	paymentConn, err := grpc.NewClient(
		paymentServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		// Cleanup: закрываем уже открытое inventoryServiceConn соединение при ошибке
		if cerr := inventoryConn.Close(); cerr != nil {
			log.Printf("failed to close inventory service grpc connection: %v", cerr)
		}
		return nil, nil, nil, fmt.Errorf("failed to create payment service grpc connection: %w", err)
	}

	inventoryClient := inventoryv1.NewInventoryServiceClient(inventoryConn)
	paymentClient := paymentv1.NewPaymentServiceClient(paymentConn)

	storage := NewOrderStorage()
	orderHandler := NewOrderHandler(storage, inventoryClient, paymentClient)

	orderServer, err := orderv1.NewServer(orderHandler)
	if err != nil {
		// Cleanup: закрываем уже открытое inventoryServiceConn соединение при ошибке
		if cerr := inventoryConn.Close(); cerr != nil {
			log.Printf("failed to close inventory service grpc connection: %v", cerr)
		}
		if cerr := paymentConn.Close(); cerr != nil {
			log.Printf("failed to close payment service grpc connection: %v", cerr)
		}
		return nil, nil, nil, fmt.Errorf("failed to create order server: %w", err)
	}

	return inventoryConn, paymentConn, orderServer, nil
}

func main() {
	inventoryConn, paymentConn, orderServer, err := initApplication()
	if err != nil {
		log.Fatalf("failed to init application: %v", err)
	}

	defer func() {
		if cerr := inventoryConn.Close(); cerr != nil {
			log.Printf("failed to close inventory service grpc connection: %v", cerr)
		}
	}()

	defer func() {
		if cerr := paymentConn.Close(); cerr != nil {
			log.Printf("failed to close payment service grpc connection: %v", cerr)
		}
	}()

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(requestTimeout))

	r.Mount("/", orderServer)

	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		log.Printf("http server listening on %s\n", server.Addr)
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("failed to start http server: %v", err)
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down http server...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("failed to shutdown http server: %v", err)
	}

	log.Println("http server stopped")
}

package v1

import (
	"github.com/qyrlabs/test-backend/inventory/internal/service"
	inventoryv1 "github.com/qyrlabs/test-backend/shared/pkg/proto/inventory/v1"
)

type api struct {
	inventoryv1.UnimplementedInventoryServiceServer

	inventoryService service.PartService
}

func NewAPI(inventoryService service.PartService) *api {
	return &api{
		inventoryService: inventoryService,
	}
}

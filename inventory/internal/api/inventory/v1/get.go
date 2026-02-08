package v1

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/qyrlabs/test-backend/inventory/internal/converter"
	"github.com/qyrlabs/test-backend/inventory/internal/model"
	inventoryv1 "github.com/qyrlabs/test-backend/shared/pkg/proto/inventory/v1"
)

// Get part info by its UUID.
func (a *api) GetPart(ctx context.Context, req *inventoryv1.GetPartRequest) (*inventoryv1.GetPartResponse, error) {
	partUUID := req.GetUuid()

	if _, err := uuid.Parse(partUUID); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid uuid format: %v", err)
	}

	part, err := a.inventoryService.Get(ctx, req.GetUuid())
	if err != nil {
		if errors.Is(err, model.ErrPartNotFound) {
			return nil, status.Errorf(codes.NotFound, "part with uuid %s is not found", req.GetUuid())
		}
		log.Printf("failed to get part with uuid %s: %v", req.GetUuid(), err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &inventoryv1.GetPartResponse{
		Part: converter.ToProtoPart(part),
	}, nil
}

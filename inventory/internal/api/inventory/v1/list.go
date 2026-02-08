package v1

import (
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/qyrlabs/test-backend/inventory/internal/converter"
	inventoryv1 "github.com/qyrlabs/test-backend/shared/pkg/proto/inventory/v1"
)

// Returns List of Parts by filter.
func (a *api) ListParts(ctx context.Context, req *inventoryv1.ListPartsRequest) (*inventoryv1.ListPartsResponse, error) {
	filteredParts, err := a.inventoryService.List(ctx, converter.ToProtoFilter(req.GetFilter()))
	if err != nil {
		log.Printf("failed to list parts: %v", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	protoParts := converter.ToProtoParts(filteredParts)

	return &inventoryv1.ListPartsResponse{
		Parts: protoParts,
	}, nil
}

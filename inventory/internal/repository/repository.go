package repository

import (
	"context"

	"github.com/qyrlabs/test-backend/inventory/internal/model"
)

type PartRepository interface {
	Get(ctx context.Context, uuid string) (*model.Part, error)
	List(ctx context.Context, filter model.PartsFilter) ([]*model.Part, error)
}

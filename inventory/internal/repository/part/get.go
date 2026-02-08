package part

import (
	"context"

	"github.com/qyrlabs/test-backend/inventory/internal/model"
	"github.com/qyrlabs/test-backend/inventory/internal/repository/converter"
)

func (r *repository) Get(ctx context.Context, uuid string) (*model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	part, ok := r.parts[uuid]
	if !ok {
		return nil, model.ErrPartNotFound
	}

	return converter.ToModelPart(part), nil
}

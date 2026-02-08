package part

import (
	"context"

	"github.com/qyrlabs/test-backend/inventory/internal/model"
)

// Returns List of Parts by filter.
func (s *service) List(ctx context.Context, filter model.PartsFilter) ([]*model.Part, error) {
	parts, err := s.partRepository.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	return parts, nil
}

package part

import (
	"context"

	"github.com/qyrlabs/test-backend/inventory/internal/model"
)

// Get part info by its UUID.
func (s *service) Get(ctx context.Context, uuid string) (*model.Part, error) {
	part, err := s.partRepository.Get(ctx, uuid)
	if err != nil {
		return nil, err
	}
	return part, nil
}

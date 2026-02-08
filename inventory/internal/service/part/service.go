package part

import (
	"github.com/qyrlabs/test-backend/inventory/internal/repository"
	def "github.com/qyrlabs/test-backend/inventory/internal/service"
)

var _ def.PartService = &service{}

type service struct {
	partRepository repository.PartRepository
}

func NewService(partRepository repository.PartRepository) *service {
	return &service{
		partRepository: partRepository,
	}
}

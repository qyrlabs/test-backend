package part

import (
	"log"
	"sync"

	def "github.com/qyrlabs/test-backend/inventory/internal/repository"
	"github.com/qyrlabs/test-backend/inventory/internal/repository/repomodel"
)

var _ def.PartRepository = &repository{}

type repository struct {
	mu    sync.RWMutex
	parts map[string]repomodel.Part
}

func NewRepository() *repository {
	repository := repository{
		parts: make(map[string]repomodel.Part),
	}
	err := repository.initParts(100)
	if err != nil {
		log.Println("failed to init parts")
	}
	return &repository
}

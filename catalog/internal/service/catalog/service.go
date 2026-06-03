package catalog

import (
	"github.com/LearLocker/streaming/catalog/internal/repository"
	def "github.com/LearLocker/streaming/catalog/internal/service"
)

var _ def.CatalogService = (*Service)(nil)

type Service struct {
	catalogRepository repository.CatalogRepository
}

func NewService(catalogRepository repository.CatalogRepository) *Service {
	return &Service{
		catalogRepository: catalogRepository,
	}
}

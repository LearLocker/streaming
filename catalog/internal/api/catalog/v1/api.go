package v1

import (
	"github.com/LearLocker/streaming/catalog/internal/service"
	catalogV1 "github.com/LearLocker/streaming/shared/pkg/proto/catalog/v1"
)

type Api struct {
	catalogV1.UnimplementedCatalogServiceServer

	catalogService service.CatalogService
}

func NewAPI(catalogService service.CatalogService) *Api {
	return &Api{
		catalogService: catalogService,
	}
}

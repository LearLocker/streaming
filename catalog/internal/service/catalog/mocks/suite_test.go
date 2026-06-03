package catalog

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/LearLocker/streaming/catalog/internal/repository/mocks"
	catalogService "github.com/LearLocker/streaming/catalog/internal/service/catalog"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context

	catalogRepository *mocks.CatalogRepository

	service *catalogService.Service
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.catalogRepository = mocks.NewCatalogRepository(s.T())

	s.service = catalogService.NewService(
		s.catalogRepository,
	)
}

func (s *ServiceSuite) TearDownTest() {
}

func TestAPIIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}

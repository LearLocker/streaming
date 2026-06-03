package v1

import (
	"context"
	"testing"

	apiV1 "github.com/LearLocker/streaming/catalog/internal/api/catalog/v1"
	"github.com/LearLocker/streaming/catalog/internal/service/mocks"
	"github.com/stretchr/testify/suite"
)

type APISuite struct {
	suite.Suite

	ctx context.Context

	catalogService *mocks.CatalogService

	api *apiV1.Api
}

func (s *APISuite) SetupTest() {
	s.ctx = context.Background()

	s.catalogService = mocks.NewCatalogService(s.T())

	s.api = apiV1.NewAPI(
		s.catalogService,
	)
}

func (s *APISuite) TearDownTest() {
}

func TestAPIIntegration(t *testing.T) {
	suite.Run(t, new(APISuite))
}

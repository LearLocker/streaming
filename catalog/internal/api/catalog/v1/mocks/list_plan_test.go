package v1

import (
	"github.com/LearLocker/streaming/catalog/internal/model"
	catalogV1 "github.com/LearLocker/streaming/shared/pkg/proto/catalog/v1"
	"github.com/brianvoe/gofakeit/v6"
)

func (s *APISuite) TestListPlanSuccess() {
	var (
		expectedPlans = []model.Plan{
			{
				PlanID:       gofakeit.UUID(),
				Name:         gofakeit.Word(),
				Price:        int64(gofakeit.Price(12, 100)),
				Currency:     gofakeit.CurrencyShort(),
				DurationDays: int32(gofakeit.RandomInt([]int{10, 20, 30, 40})),
				IsActive:     gofakeit.Bool(),
			},
			{
				PlanID:       gofakeit.UUID(),
				Name:         gofakeit.Word(),
				Price:        int64(gofakeit.Price(12, 100)),
				Currency:     gofakeit.CurrencyShort(),
				DurationDays: int32(gofakeit.RandomInt([]int{10, 20, 30, 40})),
				IsActive:     gofakeit.Bool(),
			},
		}

		req = &catalogV1.ListPlansRequest{}
	)

	s.catalogService.On("ListPlan", s.ctx, req).Return(expectedPlans, nil)

	res, err := s.api.ListPlans(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Len(res.GetPlans(), len(expectedPlans))
}

func (s *APISuite) TestListPlanServiceError() {
	var (
		serviceErr = gofakeit.Error()
		req        = &catalogV1.ListPlansRequest{}
	)

	s.catalogService.On("ListPlan", s.ctx, req).Return(nil, serviceErr)

	res, err := s.api.ListPlans(s.ctx, req)
	s.Require().Error(err)
	s.Require().ErrorIs(err, serviceErr)
	s.Require().Nil(res)
}

func (s *APISuite) TestListPlanEmptyResult() {
	var (
		req = &catalogV1.ListPlansRequest{}
	)

	s.catalogService.On("ListPlan", s.ctx, req).Return([]model.Plan{}, nil)

	res, err := s.api.ListPlans(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Empty(res.GetPlans())
}

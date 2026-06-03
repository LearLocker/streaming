package catalog

import (
	"github.com/LearLocker/streaming/catalog/internal/model"
	"github.com/brianvoe/gofakeit/v6"
)

func (s *ServiceSuite) TestListPlanSuccess() {
	var (
		active        = gofakeit.Bool()
		expectedPlans = []model.Plan{
			{
				PlanID:       gofakeit.UUID(),
				Name:         gofakeit.Word(),
				Price:        int64(gofakeit.IntRange(1, 100)),
				Currency:     gofakeit.CurrencyShort(),
				DurationDays: int32(gofakeit.IntRange(7, 365)),
				IsActive:     gofakeit.Bool(),
			},
			{
				PlanID:       gofakeit.UUID(),
				Name:         gofakeit.Word(),
				Price:        int64(gofakeit.IntRange(1, 100)),
				Currency:     gofakeit.CurrencyShort(),
				DurationDays: int32(gofakeit.IntRange(7, 365)),
				IsActive:     gofakeit.Bool(),
			},
		}
	)

	s.catalogRepository.On("ListPlans", s.ctx, active).Return(expectedPlans, nil)

	plans, err := s.catalogRepository.ListPlans(s.ctx, active)
	s.Require().NoError(err)
	s.Require().NotNil(plans)
	s.Require().Len(plans, len(expectedPlans))
	s.Require().Equal(expectedPlans, plans)
}

func (s *ServiceSuite) TestListPlanRepoError() {
	var (
		repoErr = gofakeit.Error()
		active  = gofakeit.Bool()
	)

	s.catalogRepository.On("ListPlans", s.ctx, active).Return(nil, repoErr)

	plans, err := s.catalogRepository.ListPlans(s.ctx, active)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
	s.Require().Nil(plans)
}

func (s *ServiceSuite) TestListPlanEmptyResult() {
	var (
		active = gofakeit.Bool()
	)

	s.catalogRepository.On("ListPlans", s.ctx, active).Return([]model.Plan{}, nil)

	plans, err := s.catalogRepository.ListPlans(s.ctx, active)
	s.Require().NoError(err)
	s.Require().Empty(plans)
}

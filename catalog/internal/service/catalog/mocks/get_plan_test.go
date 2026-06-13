package catalog

import (
	"github.com/LearLocker/streaming/catalog/internal/model"
	"github.com/brianvoe/gofakeit/v6"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServiceSuite) TestGetPlanSuccess() {
	var (
		planId       = gofakeit.UUID()
		planName     = gofakeit.Word()
		planPrice    = gofakeit.Price(12, 100)
		planCurrency = gofakeit.CurrencyShort()
		durationDays = gofakeit.RandomInt([]int{10, 20, 30, 40})
		planIsActive = gofakeit.Bool()

		modelPlan = model.Plan{
			PlanID:       planId,
			Name:         planName,
			Price:        int64(planPrice),
			Currency:     planCurrency,
			DurationDays: int32(durationDays),
			IsActive:     planIsActive,
		}
	)

	s.catalogRepository.On("GetPlan", s.ctx, planId).Return(modelPlan, nil)

	res, err := s.service.GetPlan(s.ctx, planId)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Equal(modelPlan, res)

}

func (s *ServiceSuite) TestGetPlanNotFound() {
	var (
		planId = gofakeit.UUID()
	)

	s.catalogRepository.On("GetPlan", s.ctx, planId).Return(nil, model.ErrPlanNotFound)

	res, err := s.service.GetPlan(s.ctx, planId)
	s.Require().Error(err)
	s.Require().Nil(res)

	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Require().Equal(codes.NotFound, st.Code())
}

func (s *ServiceSuite) TestGetPlanRepoError() {
	var (
		planId  = gofakeit.UUID()
		repoErr = gofakeit.Error()
	)

	s.catalogRepository.On("GetPlan", s.ctx, planId).Return(nil, repoErr)

	res, err := s.service.GetPlan(s.ctx, planId)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
	s.Require().Nil(res)
}

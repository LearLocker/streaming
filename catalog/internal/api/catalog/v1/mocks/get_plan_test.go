package v1

import (
	"github.com/LearLocker/streaming/catalog/internal/converter"
	"github.com/LearLocker/streaming/catalog/internal/model"
	catalogV1 "github.com/LearLocker/streaming/shared/pkg/proto/catalog/v1"
	"github.com/brianvoe/gofakeit/v6"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *APISuite) TestGetPlanSuccess() {
	var (
		planId       = gofakeit.UUID()
		planName     = gofakeit.Word()
		planPrice    = gofakeit.Price(12, 100)
		planCurrency = gofakeit.CurrencyShort()
		durationDays = gofakeit.RandomInt([]int{10, 20, 30, 40})
		planIsActive = gofakeit.Bool()

		req = &catalogV1.GetPlanRequest{
			PlanId: planId,
		}

		modelPlan = &model.Plan{
			PlanID:       planId,
			Name:         planName,
			Price:        int64(planPrice),
			Currency:     planCurrency,
			DurationDays: int32(durationDays),
			IsActive:     planIsActive,
		}

		expectedProtoPlan = converter.PlanToProto(modelPlan)
	)

	s.catalogService.On("GetPlan", s.ctx, planId).Return(modelPlan, nil)

	res, err := s.api.GetPlan(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Equal(expectedProtoPlan.PlanId, res.GetPlan().GetPlanId())
	s.Require().Equal(expectedProtoPlan.Name, res.GetPlan().GetName())
	s.Require().Equal(expectedProtoPlan.Price, res.GetPlan().GetPrice())

}

func (s *APISuite) TestGetPlanNotFound() {
	var (
		planId = gofakeit.UUID()

		req = &catalogV1.GetPlanRequest{
			PlanId: planId,
		}
	)

	s.catalogService.On("GetPlan", s.ctx, planId).Return(model.Plan{}, model.ErrPlanNotFound)

	res, err := s.api.GetPlan(s.ctx, req)
	s.Require().Error(err)
	s.Require().Nil(res)

	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Require().Equal(codes.NotFound, st.Code())
}

func (s *APISuite) TestGetPlanServiceError() {
	var (
		planId     = gofakeit.UUID()
		serviceErr = gofakeit.Error()

		req = &catalogV1.GetPlanRequest{
			PlanId: planId,
		}
	)

	s.catalogService.On("GetPlan", s.ctx, planId).Return(model.Plan{}, serviceErr)

	res, err := s.api.GetPlan(s.ctx, req)
	s.Require().Error(err)
	s.Require().ErrorIs(err, serviceErr)
	s.Require().Nil(res)
}

package repository

import (
	"context"

	userProto "github.com/goplaceapp/goplace-user/api/v1"
	"github.com/goplaceapp/goplace-user/internal/tenant/domain"
)

func (r *TenantRepository) RequestDemo(ctx context.Context, req *userProto.DemoRequest) (*userProto.DemoResponse, error) {
	requestedDemo := &domain.Demo{
		Name:           req.GetParams().GetName(),
		Email:          req.GetParams().GetEmail(),
		PhoneNumber:    req.GetParams().GetPhoneNumber(),
		Country:        req.GetParams().GetCountry(),
		RestaurantName: req.GetParams().GetRestaurantName(),
		BranchesNo:     req.GetParams().GetBranchesNo(),
		FirstTimeCrm:   req.GetParams().GetFirstTimeCrm(),
		SystemName:     req.GetParams().GetSystemName(),
	}

	result := r.SharedDbConnection.Table("demo_requests").Create(&requestedDemo)
	if result.Error != nil {
		return nil, result.Error
	}

	return &userProto.DemoResponse{
		Result: "Demo requested successfully",
	}, nil
}

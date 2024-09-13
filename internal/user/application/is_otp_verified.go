package application

import (
	"context"

	userProto "github.com/goplaceapp/goplace-user/api/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *UserService) IsOtpVerified(ctx context.Context, req *emptypb.Empty) (*userProto.IsOtpVerifiedResponse, error) {
	res, err := s.Repository.IsOtpVerified(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

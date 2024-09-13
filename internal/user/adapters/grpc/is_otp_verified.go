package grpc

import (
	"context"

	userProto "github.com/goplaceapp/goplace-user/api/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *UserServiceServer) IsOtpVerified(ctx context.Context, req *emptypb.Empty) (*userProto.IsOtpVerifiedResponse, error) {
	res, err := s.userService.IsOtpVerified(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

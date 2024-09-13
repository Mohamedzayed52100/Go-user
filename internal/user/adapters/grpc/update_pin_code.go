package grpc

import (
	"context"

	userProto "github.com/goplaceapp/goplace-user/api/v1"
)

func (s *UserServiceServer) UpdatePinCode(ctx context.Context, req *userProto.UpdatePinCodeRequest) (*userProto.UpdatePinCodeResponse, error) {
	res, err := s.userService.UpdatePinCode(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

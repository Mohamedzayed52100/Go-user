package application

import (
	"context"

	userProto "github.com/goplaceapp/goplace-user/api/v1"
)

func (s *UserService) CheckPinCode(ctx context.Context, req *userProto.PinCodeRequest) (*userProto.PinCodeResponse, error) {
	res, err := s.Repository.CheckPinCode(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

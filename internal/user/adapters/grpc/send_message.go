package grpc

import (
	"context"

	userProto "github.com/goplaceapp/goplace-user/api/v1"
)

func (s *UserServiceServer) SendMessage(ctx context.Context, req *userProto.SendMessageRequest) (*userProto.SendMessageResponse, error) {
	res, err := s.userService.SendMessage(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

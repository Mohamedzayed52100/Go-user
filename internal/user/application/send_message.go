package application

import (
	"context"

	userProto "github.com/goplaceapp/goplace-user/api/v1"
)

func (s *UserService) SendMessage(ctx context.Context, req *userProto.SendMessageRequest) (*userProto.SendMessageResponse, error) {
	res, err := s.Repository.SendMessage(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

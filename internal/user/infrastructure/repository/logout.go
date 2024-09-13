package repository

import (
	"context"
	"net/http"
	"time"

	"github.com/goplaceapp/goplace-common/pkg/auth"
	"github.com/goplaceapp/goplace-common/pkg/meta"
	externalUserDomain "github.com/goplaceapp/goplace-user/pkg/userservice/domain"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (r *UserRepository) Logout(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	sessionId := auth.GetSessionIDFromToken(ctx.Value(meta.AuthorizationContextKey.String()).(string))

	if err := r.SharedDbConnection.Model(&externalUserDomain.UserSession{}).Where("session_id = ?", sessionId).Update("logout_time", time.Now().UTC()).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	return &emptypb.Empty{}, nil
}

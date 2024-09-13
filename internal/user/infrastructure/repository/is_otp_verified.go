package repository

import (
	"context"

	"github.com/goplaceapp/goplace-common/pkg/auth"
	"github.com/goplaceapp/goplace-common/pkg/meta"
	userProto "github.com/goplaceapp/goplace-user/api/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (r *UserRepository) IsOtpVerified(ctx context.Context, _ *emptypb.Empty) (*userProto.IsOtpVerifiedResponse, error) {
	token := ctx.Value(meta.AuthorizationContextKey.String()).(string)
	email := auth.GetAccountMetadataFromToken(token).Email

	isVerified := false
	r.GetSharedDB().
		Table("users").
		Joins("JOIN user_otps ON user_otps.user_id = users.id").Where("UPPER(user_otps.type) = 'LOGIN' AND users.email = ? AND verified = true", email).Select("user_otps.verified").Scan(&isVerified)

	return &userProto.IsOtpVerifiedResponse{Result: isVerified}, nil
}

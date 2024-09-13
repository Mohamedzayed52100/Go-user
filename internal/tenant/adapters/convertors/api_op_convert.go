package convertors

import (
	userProto "github.com/goplaceapp/goplace-user/api/v1"
	"github.com/goplaceapp/goplace-user/pkg/tenantservice/domain"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func BuildTenantResponse(proto *domain.Tenant) *userProto.Tenant {
	return &userProto.Tenant{
		Id:        proto.ID,
		Domain:    proto.Domain,
		DbName:    proto.DbName,
		CreatedAt: timestamppb.New(proto.CreatedAt),
		UpdatedAt: timestamppb.New(proto.UpdatedAt),
	}
}

package common

import (
	"math"

	userProto "github.com/goplaceapp/goplace-user/api/v1"
)

func BuildPaginationResponse(pagination *userProto.UPaginationParams, length int32, offset int32) *userProto.UPagination {
	var lastPage int32
	if pagination.GetPerPage() != 0 {
		lastPage = int32(math.Ceil(float64(length) / float64(pagination.GetPerPage())))
	} else {
		lastPage = 0
	}

	return &userProto.UPagination{
		Total:       length,
		PerPage:     pagination.GetPerPage(),
		CurrentPage: pagination.GetCurrentPage(),
		LastPage:    lastPage,
		From:        offset + 1,
		To:          offset + pagination.GetPerPage(),
	}
}

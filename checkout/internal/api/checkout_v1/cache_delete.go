package checkout_v1

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
	"route256/checkout/internal/models"
	desc "route256/checkout/pkg/checkout_v1"
)

func (i *Implementation) CacheDelete(ctx context.Context, req *desc.CacheDeleteRequest) (*emptypb.Empty, error) {
	var keys models.Keys
	keys = append(keys, req.GetKeys()...)
	err := i.BusinessLogic.CacheDelete(ctx, keys)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

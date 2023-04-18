package checkout_v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/emptypb"
	"route256/checkout/internal/models"
	desc "route256/checkout/pkg/checkout_v1"
	"route256/libs/grpcresponse"
)

func (i *Implementation) DeleteFromCart(ctx context.Context, req *desc.DeleteFromCartRequest) (*emptypb.Empty, error) {
	err := i.BusinessLogic.DeleteFromCart(ctx, models.UserID(req.GetUser()),
		models.Sku(req.GetSku()), models.Count(req.GetCount()))

	if err != nil {
		return nil, grpcresponse.Error(err, codes.Internal, "deleting from cart")
	}
	return &emptypb.Empty{}, nil
}

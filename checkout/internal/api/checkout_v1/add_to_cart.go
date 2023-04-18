package checkout_v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/emptypb"
	"route256/checkout/internal/models"
	desc "route256/checkout/pkg/checkout_v1"
	"route256/libs/grpcresponse"
)

func (i *Implementation) AddToCart(ctx context.Context, req *desc.AddToCartRequest) (*emptypb.Empty, error) {
	err := i.BusinessLogic.AddToCart(ctx, models.UserID(req.GetUser()),
		models.Sku(req.GetSku()), models.Count(req.GetCount()))

	if err != nil {
		return nil, grpcresponse.Error(err, codes.Internal, "adding to cart")
	}
	return &emptypb.Empty{}, nil
}

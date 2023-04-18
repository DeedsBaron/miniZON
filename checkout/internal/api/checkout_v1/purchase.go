package checkout_v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"route256/checkout/internal/models"
	desc "route256/checkout/pkg/checkout_v1"
	"route256/libs/grpcresponse"
)

func (i *Implementation) Purchase(ctx context.Context, req *desc.PurchaseRequest) (*desc.PurchaseResponse, error) {
	orderID, err := i.BusinessLogic.Purchase(ctx, models.UserID(req.GetUser()))
	if err != nil {
		return nil, grpcresponse.Error(err, codes.Internal, "purchasing")
	}

	return &desc.PurchaseResponse{
		OrderId: int64(orderID),
	}, nil
}

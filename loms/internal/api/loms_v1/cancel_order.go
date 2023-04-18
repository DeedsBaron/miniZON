package loms_v1

import (
	"context"
	"route256/loms/internal/models"
	desc "route256/loms/pkg/loms_v1"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) CancelOrder(ctx context.Context, req *desc.CancelOrderRequest) (*emptypb.Empty, error) {
	err := i.BusinessLogic.CancelOrder(ctx, models.OrderID(req.GetOrderId()))
	if err != nil {
		return nil, errors.WithMessage(err, "canceling order")
	}
	return &emptypb.Empty{}, nil
}

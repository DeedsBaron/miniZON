package loms_v1

import (
	"context"
	"route256/libs/grpcresponse"
	"route256/loms/internal/models"
	desc "route256/loms/pkg/loms_v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) OrderPayed(ctx context.Context, req *desc.OrderPayedRequest) (*emptypb.Empty, error) {
	err := i.BusinessLogic.OrderPayed(ctx, models.OrderID(req.GetOrderId()))
	if err != nil {
		return nil, grpcresponse.Error(err, codes.Internal, "paying order")
	}
	return &emptypb.Empty{}, nil
}

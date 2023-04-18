package loms_v1

import (
	"context"
	"route256/libs/grpcresponse"
	"route256/loms/internal/models"
	desc "route256/loms/pkg/loms_v1"

	"google.golang.org/grpc/codes"
)

func (i *Implementation) ListOrder(ctx context.Context, req *desc.ListOrderRequest) (*desc.ListOrderResponse, error) {
	orderList, err := i.BusinessLogic.ListOrder(ctx, models.OrderID(req.GetOrderId()))
	if err != nil {
		return nil, grpcresponse.Error(err, codes.Internal, "listing order")
	}
	resp := &desc.ListOrderResponse{
		User:  int64(orderList.User),
		Items: make([]*desc.Item, 0, len(orderList.Items)),
	}
	switch orderList.Status {
	case models.StatusNew:
		resp.Status = desc.Status_new
	case models.StatusAwaitingPayment:
		resp.Status = desc.Status_awaiting_payment
	case models.StatusPayed:
		resp.Status = desc.Status_payed
	case models.StatusFailed:
		resp.Status = desc.Status_failed
	case models.StatusCanceled:
		resp.Status = desc.Status_canceled
	}
	for _, item := range orderList.Items {
		resp.Items = append(resp.Items, &desc.Item{
			Sku:   item.Sku,
			Count: item.Count,
		})
	}
	return resp, nil
}

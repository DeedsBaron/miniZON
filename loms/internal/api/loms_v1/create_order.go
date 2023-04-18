package loms_v1

import (
	"context"
	"route256/libs/grpcresponse"
	"route256/loms/internal/models"
	desc "route256/loms/pkg/loms_v1"

	"google.golang.org/grpc/codes"
)

func (i *Implementation) CreateOrder(ctx context.Context, req *desc.CreateOrderRequest) (*desc.CreateOrderResponse, error) {
	createOrder := models.CreateOrder{
		User:  req.GetUser(),
		Items: make([]models.OrderItem, 0, len(req.GetItems())),
	}
	reqItems := req.GetItems()
	for _, item := range reqItems {
		createOrder.Items = append(createOrder.Items, models.OrderItem{
			Sku:   item.Sku,
			Count: item.Count,
		})
	}

	orderID, err := i.BusinessLogic.CreateOrder(ctx, createOrder)
	if err != nil {
		return nil, grpcresponse.Error(err, codes.Internal, "creating order")
	}
	return &desc.CreateOrderResponse{
		OrderId: int64(orderID),
	}, nil
}

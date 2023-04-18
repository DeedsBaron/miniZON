package checkout_v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"route256/checkout/internal/models"
	desc "route256/checkout/pkg/checkout_v1"
	"route256/libs/grpcresponse"
)

func (i *Implementation) ListCart(ctx context.Context, req *desc.ListCartRequest) (*desc.ListCartResponse, error) {
	cartList, err := i.BusinessLogic.ListCart(ctx, models.UserID(req.GetUser()))
	if err != nil {
		return nil, grpcresponse.Error(err, codes.Internal, "listing cart")
	}
	resp := &desc.ListCartResponse{
		Items:      make([]*desc.Item, 0, len(cartList.Items)),
		TotalPrice: cartList.TotalPrice,
	}
	for _, item := range cartList.Items {
		resp.Items = append(resp.Items, &desc.Item{
			Sku:   uint32(item.Sku),
			Count: uint32(item.Count),
			Name:  item.Name,
			Price: item.Price,
		})
	}
	resp.TotalPrice = cartList.TotalPrice
	return resp, nil
}

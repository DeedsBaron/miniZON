package loms

import (
	"context"

	"google.golang.org/grpc"
	"route256/checkout/internal/config"
	"route256/checkout/internal/models"
	"route256/libs/clientwrapper"
	lomsServiceAPI "route256/loms/pkg/loms_v1"
)

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate minimock -i LomsService -o ./mocks/ -s "_minimock.go"

type LomsService interface {
	Stocks(ctx context.Context, sku uint32) ([]models.Stock, error)
	CreateOrder(ctx context.Context, userID models.UserID, items []models.Item) (models.OrderID, error)
}

type client struct {
	lomsClient lomsServiceAPI.LomsV1Client
	Conn       *grpc.ClientConn
}

func NewClient(ctx context.Context) *client {
	conn := clientwrapper.NewGrpcConnection(ctx, config.Data.Services.Loms)
	return &client{
		lomsClient: lomsServiceAPI.NewLomsV1Client(conn),
		Conn:       conn,
	}
}

func (c *client) Stocks(ctx context.Context, sku uint32) ([]models.Stock, error) {
	resp, err := c.lomsClient.Stocks(ctx, &lomsServiceAPI.StocksRequest{
		Sku: sku,
	})
	if err != nil {
		return nil, err
	}
	ds := make([]models.Stock, 0, len(resp.Stocks))
	for _, stock := range resp.Stocks {
		ds = append(ds, models.Stock{
			WarehouseID: stock.GetWarehouseId(),
			Count:       stock.GetCount(),
		})
	}
	return ds, nil
}

func (c *client) CreateOrder(ctx context.Context, userID models.UserID, items []models.Item) (models.OrderID, error) {
	req := &lomsServiceAPI.CreateOrderRequest{
		User:  int64(userID),
		Items: make([]*lomsServiceAPI.Item, 0, len(items)),
	}

	for _, item := range items {
		req.Items = append(req.Items, &lomsServiceAPI.Item{
			Sku:   uint32(item.Sku),
			Count: uint32(item.Count),
		})
	}

	resp, err := c.lomsClient.CreateOrder(ctx, req)
	if err != nil {
		return 0, err
	}

	orderID := models.OrderID(resp.OrderId)
	return orderID, nil
}

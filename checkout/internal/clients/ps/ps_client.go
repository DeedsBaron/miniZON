package ps

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"route256/checkout/internal/config"
	"route256/checkout/internal/models"
	psServiceAPI "route256/checkout/internal/pb/ps"
	"route256/libs/clientwrapper"
)

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate minimock -i ProductService -o ./mocks/ -s "_minimock.go"

type ProductService interface {
	GetProduct(ctx context.Context, sku uint32) (*models.ProductInfo, error)
}

type client struct {
	ps   psServiceAPI.ProductServiceClient
	Conn *grpc.ClientConn
}

func NewClient(ctx context.Context) *client {
	conn := clientwrapper.NewGrpcConnection(ctx, config.Data.Services.Ps)
	return &client{
		ps:   psServiceAPI.NewProductServiceClient(conn),
		Conn: conn,
	}
}

func (c *client) GetProduct(ctx context.Context, sku uint32) (*models.ProductInfo, error) {

	product, err := c.ps.GetProduct(ctx, &psServiceAPI.GetProductRequest{
		Token: config.Data.Token,
		Sku:   sku,
	})
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("ps client GetProduct for sku: %d", sku))
	}

	productInfo := &models.ProductInfo{
		Name:  product.GetName(),
		Price: product.GetPrice(),
	}

	return productInfo, nil
}

func (c *client) ListSkus(ctx context.Context, startAftersku, count uint32) (*[]models.Sku, error) {
	resp, err := c.ps.ListSkus(ctx, &psServiceAPI.ListSkusRequest{
		Token:         config.Data.Token,
		StartAfterSku: startAftersku,
		Count:         count,
	})
	if err != nil {
		return nil, errors.Wrap(err, "ps client ListSkus")
	}
	skuList := make([]models.Sku, 0, len(resp.GetSkus()))
	for i, _ := range skuList {
		skuList[i] = models.Sku(resp.GetSkus()[i])
	}

	return &skuList, nil
}

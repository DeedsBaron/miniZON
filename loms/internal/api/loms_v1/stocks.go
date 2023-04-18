package loms_v1

import (
	"context"
	"route256/libs/grpcresponse"
	desc "route256/loms/pkg/loms_v1"

	"google.golang.org/grpc/codes"
)

type StocksItem struct {
	WarehouseID int64
	Count       uint64
}

type Stocks struct {
	Stocks []StocksItem
}

func (i *Implementation) Stocks(ctx context.Context, req *desc.StocksRequest) (*desc.StocksResponse, error) {
	stocks, err := i.BusinessLogic.GetStocks(ctx, req.GetSku())
	if err != nil {
		return nil, grpcresponse.Error(err, codes.Internal, "getting stocks")
	}

	stocksResp := &desc.StocksResponse{
		Stocks: make([]*desc.Stock, 0, len(stocks.Stocks)),
	}
	for _, stockItem := range stocks.Stocks {
		stocksResp.Stocks = append(stocksResp.Stocks, &desc.Stock{
			WarehouseId: int64(stockItem.WarehouseID),
			Count:       stockItem.Count,
		})
	}

	return stocksResp, nil
}

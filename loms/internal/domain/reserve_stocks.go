package domain

import (
	"context"
	"errors"

	"route256/loms/internal/models"
)

var (
	ErrNotEnoughItemsOnStocks = errors.New("not enough items on stocks")
)

func (b *Domain) ReserveStocks(ctxTX context.Context, orderID models.OrderID, itemSku uint32, count uint32) error {
	var needToReserve = make(models.ItemStocks, 1)

	stocks, err := b.repo.GetStocks(ctxTX, itemSku)
	if err != nil {
		return err
	}

	var (
		reservedCount uint64
	)

	for _, stockItem := range stocks.Stocks {
		left := uint64(count) - reservedCount
		if left == 0 {
			break
		}
		if stockItem.Count >= left {
			needToReserve[stockItem.WarehouseID] = left
			reservedCount += left
		} else {
			needToReserve[stockItem.WarehouseID] = stockItem.Count
			reservedCount += stockItem.Count
		}
	}

	if reservedCount != uint64(count) {
		return ErrNotEnoughItemsOnStocks
	}
	for warehouseID, needToReserveCount := range needToReserve {
		if err = b.repo.ReserveStocks(ctxTX, orderID, itemSku, needToReserveCount, warehouseID); err != nil {
			return err
		}
	}

	return nil
}

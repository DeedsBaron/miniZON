package repository

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"route256/libs/logger"
	"route256/loms/internal/models"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
)

func (r *LomsRepository) ReserveStocks(ctx context.Context, orderId models.OrderID, itemSku uint32, needToReserveCount uint64, warehouseId models.WarehouseID) error {
	db := r.GetQueryEngine(ctx)
	childSpan, _ := opentracing.StartSpanFromContext(ctx, "repo: reserve_stocks")
	defer childSpan.Finish()
	sq := squirrel.Expr(
		`insert into reservation (order_id, item_sku, warehouse_id, reserved_count)
			values($1, $2, $3, $4);`,
		orderId, itemSku, warehouseId, needToReserveCount)

	rawQuery, args, err := sq.ToSql()
	if err != nil {
		logger.Errorf("tosql error %v", err)
		return models.ErrInternal
	}
	var id []int
	if err = pgxscan.Select(ctx, db, &id, rawQuery, args...); err != nil {
		logger.Errorf("reserve stocks: orderId=%d, itemSku=%d: database error: %v", orderId, itemSku, err)
		return models.ErrInternal
	}
	return nil
}

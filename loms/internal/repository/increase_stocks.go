package repository

import (
	"context"

	"route256/libs/logger"
	"route256/loms/internal/models"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
)

func (r *LomsRepository) IncreaseStock(ctx context.Context, orderID models.OrderID) error {
	db := r.GetQueryEngine(ctx)

	sq := squirrel.Expr(
		`update stocks
		set count = count + r.reserved_count
		from (
			select item_sku, warehouse_id, reserved_count from reservation where order_id = $1			
		) as r
		where stocks.warehouse_id = r.warehouse_id and stocks.item_sku = r.item_sku`,
		orderID)

	rawQuery, args, err := sq.ToSql()
	if err != nil {
		logger.Errorf("tosql error %v", err)
		return models.ErrInternal
	}
	var id []int
	if err = pgxscan.Select(ctx, db, &id, rawQuery, args...); err != nil {
		logger.Errorf("increase stock failed for orderId=%d, database error: %v", orderID, err)
		return models.ErrInternal
	}

	return nil
}

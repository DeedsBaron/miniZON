package repository

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"route256/libs/logger"
	sq "route256/libs/squirrel"
	"route256/loms/internal/models"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
)

func (r *LomsRepository) UpdateOrderStatus(ctx context.Context, orderIDs []models.OrderID, status models.Status) error {
	db := r.GetQueryEngine(ctx)
	childSpan, _ := opentracing.StartSpanFromContext(ctx, "repo: update_order_status")
	defer childSpan.Finish()
	sq := sq.PgSb().Update("orders").
		Set("status", status).
		Where(squirrel.Eq{"id": orderIDs})

	rawQuery, args, err := sq.ToSql()
	if err != nil {
		logger.Errorf("tosql error %v", err)
		return models.ErrInternal
	}
	var id []int
	if err = pgxscan.Select(ctx, db, &id, rawQuery, args...); err != nil {
		logger.Errorf("update order status=%s failed for orderIds=%v, database error: %v", status, orderIDs, err)
		return models.ErrInternal
	}
	return nil
}

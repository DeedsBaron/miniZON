package repository

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"route256/libs/logger"
	"route256/loms/internal/models"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
)

func (r *LomsRepository) GetOrderStatus(ctx context.Context, orderID models.OrderID) (models.Status, error) {
	db := r.GetQueryEngine(ctx)
	childSpan, _ := opentracing.StartSpanFromContext(ctx, "repo: get_order_status")
	defer childSpan.Finish()
	sq := squirrel.Expr(
		`select status
			from orders
			where id = $1;`,
		orderID)

	rawQuery, args, err := sq.ToSql()
	if err != nil {
		logger.Errorf("tosql error %v", err)
		return "", models.ErrInternal
	}
	var status []models.Status
	if err = pgxscan.Select(ctx, db, &status, rawQuery, args...); err != nil {
		logger.Errorf("getting order status failed for orderId=%d, database error: %v", orderID, err)
		return "", models.ErrInternal
	}
	if status == nil {
		return "", models.ErrNotFound
	}
	return status[0], nil
}

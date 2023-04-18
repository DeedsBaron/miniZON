package repository

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"route256/libs/logger"
	"route256/loms/internal/models"
	"route256/loms/internal/repository/schema"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
)

func (r *LomsRepository) ListOrder(ctx context.Context, orderID models.OrderID) (*models.Order, error) {
	db := r.GetQueryEngine(ctx)
	childSpan, _ := opentracing.StartSpanFromContext(ctx, "repo: list_order")
	defer childSpan.Finish()
	sq := squirrel.Expr(
		`select item_sku, SUM(reserved_count) as reserved_count
		from reservation
		where order_id = $1
		group by (item_sku);`,
		orderID)

	rawQuery, args, err := sq.ToSql()
	if err != nil {
		logger.Errorf("tosql error %v", err)
		return nil, models.ErrInternal
	}
	var lo []schema.ListOrder
	if err = pgxscan.Select(ctx, db, &lo, rawQuery, args...); err != nil {
		logger.Errorf("listing order %d: database error: %v\n", orderID, err)
		return nil, models.ErrInternal
	}
	if lo == nil {
		return nil, models.ErrNotFound
	}
	status, err := r.GetOrderStatus(ctx, orderID)
	if err != nil {
		return nil, models.ErrInternal
	}
	user, err := r.GetOrderUser(ctx, orderID)
	if err != nil {
		return nil, models.ErrInternal
	}

	orderList := &models.Order{
		Status: status,
		User:   user,
		Items:  make([]models.OrderItem, 0, len(lo)),
	}
	for _, item := range lo {
		orderList.Items = append(orderList.Items, models.OrderItem{
			Sku:   item.ItemSku,
			Count: item.ReservedCount,
		})
	}

	return orderList, nil
}

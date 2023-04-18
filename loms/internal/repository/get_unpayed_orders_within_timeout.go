package repository

import (
	"context"
	"time"

	"github.com/opentracing/opentracing-go"
	"route256/libs/logger"
	"route256/loms/internal/config"
	"route256/loms/internal/models"
	"route256/loms/internal/repository/schema"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
)

func (r *LomsRepository) GetUnpayedOrdersWithinTimeout(ctx context.Context) ([]models.OrderID, error) {
	db := r.GetQueryEngine(ctx)
	childSpan, _ := opentracing.StartSpanFromContext(ctx, "repo: get_unpaided_orders_within_timeout")
	defer childSpan.Finish()
	sq := squirrel.Expr(
		`select id 
			from orders 
			where status = 'awaiting_payment' and
			created_at < ((SELECT $1 AT TIME ZONE 'Europe/Moscow') - $2::INTERVAL);`,
		time.Now(), config.Data.CronJobs.CancelReservationDueTimeoutJob.OrderToBePayedTimeout.String())

	rawQuery, args, err := sq.ToSql()
	if err != nil {
		logger.Errorf("tosql error %v", err)
		return nil, models.ErrInternal
	}
	var or []schema.Order
	if err = pgxscan.Select(ctx, db, &or, rawQuery, args...); err != nil {
		logger.Errorf("cheking unpayed orders database error, err:%v\n", err)
		return nil, models.ErrInternal
	}

	orders := make([]models.OrderID, 0, len(or))
	for _, order := range or {
		orders = append(orders, models.OrderID(order.Id))
	}

	return orders, nil
}

package repository

import (
	"context"
	"time"

	"github.com/opentracing/opentracing-go"
	"route256/libs/logger"
	"route256/loms/internal/models"
	"route256/loms/internal/repository/schema"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
)

func (r *LomsRepository) CreateOrder(ctx context.Context, createOrder models.CreateOrder) (models.OrderID, error) {
	db := r.GetQueryEngine(ctx)
	childSpan, _ := opentracing.StartSpanFromContext(ctx, "repo: create_order")
	defer childSpan.Finish()

	sq := squirrel.Expr(
		`insert into orders (user_id, status, created_at)
			values($1, $2, $3)
			returning "id";`,
		createOrder.User, models.StatusNew, time.Now())

	rawQuery, args, err := sq.ToSql()
	if err != nil {
		logger.Errorf("tosql error %v", err)
		return 0, models.ErrInternal
	}
	var or []schema.Order
	if err = pgxscan.Select(ctx, db, &or, rawQuery, args...); err != nil {
		logger.Errorf("create order %v: database error: %v\n", createOrder, err)
		return 0, models.ErrInternal
	}
	orderID := models.OrderID(or[0].Id)

	return orderID, nil
}

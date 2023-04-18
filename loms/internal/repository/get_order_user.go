package repository

import (
	"context"
	"log"

	"github.com/opentracing/opentracing-go"
	"route256/libs/logger"
	"route256/loms/internal/models"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
)

func (r *LomsRepository) GetOrderUser(ctx context.Context, orderID models.OrderID) (models.User, error) {
	db := r.GetQueryEngine(ctx)
	childSpan, _ := opentracing.StartSpanFromContext(ctx, "repo: get_order_user")
	defer childSpan.Finish()
	sq := squirrel.Expr(
		`select user_id
			from orders
			where id = $1;`,
		orderID)

	rawQuery, args, err := sq.ToSql()
	if err != nil {
		logger.Errorf("tosql error %v", err)
		return 0, models.ErrInternal
	}
	var user []models.User
	if err = pgxscan.Select(ctx, db, &user, rawQuery, args...); err != nil {
		log.Printf("getting order status failed for orderId=%d, database error: %v", orderID, err)
		return 0, models.ErrInternal
	}
	if user == nil {
		return 0, models.ErrNotFound
	}
	return user[0], nil
}

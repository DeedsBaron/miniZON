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

func (r *LomsRepository) SetOutbox(ctx context.Context, orderID models.OrderID, oldStatus, newStatus models.Status) error {
	db := r.GetQueryEngine(ctx)
	childSpan, _ := opentracing.StartSpanFromContext(ctx, "repo: set_outbox")
	defer childSpan.Finish()

	sq := squirrel.Expr(
		`insert into outbox (order_id, old_status, new_status, changed_at)
			values($1, $2, $3, $4)`,
		orderID, oldStatus, newStatus, time.Now())

	rawQuery, args, err := sq.ToSql()
	if err != nil {
		logger.Errorf("tosql error %v", err)
		return models.ErrInternal
	}
	var id []schema.Id
	if err = pgxscan.Select(ctx, db, &id, rawQuery, args...); err != nil {
		logger.Errorf("set outbox orderID %v: database error: %v\n", orderID, err)
		return models.ErrInternal
	}

	return nil
}

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

func (r *LomsRepository) UpdateOutbox(ctx context.Context, ids []models.OutboxID) error {
	db := r.GetQueryEngine(ctx)
	childSpan, _ := opentracing.StartSpanFromContext(ctx, "repo: update_outbox")
	defer childSpan.Finish()
	sq := sq.PgSb().Update("outbox").
		Set("is_sent", "sent").
		Where(squirrel.Eq{"id": ids})

	rawQuery, args, err := sq.ToSql()
	if err != nil {
		logger.Errorf("tosql error %v", err)
		return models.ErrInternal
	}
	var id []int
	if err = pgxscan.Select(ctx, db, &id, rawQuery, args...); err != nil {
		logger.Errorf("update outbox failed for ids=%v, database error: %v", ids, err)
		return models.ErrInternal
	}
	return nil
}

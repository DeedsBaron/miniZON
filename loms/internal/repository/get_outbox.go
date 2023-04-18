package repository

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"route256/libs/logger"
	"route256/loms/internal/config"
	"route256/loms/internal/models"
	"route256/loms/internal/repository/schema"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
)

func (r *LomsRepository) GetOutbox(ctx context.Context) ([]models.Outbox, error) {
	db := r.GetQueryEngine(ctx)
	childSpan, _ := opentracing.StartSpanFromContext(ctx, "repo: get_outbox")
	defer childSpan.Finish()
	sq := squirrel.Expr(
		`select id, order_id, old_status, new_status, changed_at
			from outbox
			where is_sent = 'pending'
			ORDER BY changed_at ASC
			limit $1`, config.Data.CronJobs.ReadOutBoxSendJob.BatchSizeToRead)

	rawQuery, args, err := sq.ToSql()
	if err != nil {
		logger.Errorf("tosql error %v", err)
		return nil, models.ErrInternal
	}
	var o []schema.Outbox
	if err = pgxscan.Select(ctx, db, &o, rawQuery, args...); err != nil {
		logger.Errorf("getting first %d msgs of outbox failed, database error: %v",
			config.Data.CronJobs.ReadOutBoxSendJob.BatchSizeToRead, err)
		return nil, models.ErrInternal
	}
	outbox := make([]models.Outbox, 0, len(o))
	for _, row := range o {
		outbox = append(outbox, models.Outbox{
			ID:        models.OutboxID(row.ID),
			OrderID:   models.OrderID(row.OrderID),
			OldStatus: models.Status(row.OldStatus),
			NewStatus: models.Status(row.NewStatus),
			ChangedAt: row.ChangedAt,
		})
	}

	return outbox, nil
}

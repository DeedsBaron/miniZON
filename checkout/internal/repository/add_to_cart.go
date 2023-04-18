package repository

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/opentracing/opentracing-go"
	"route256/checkout/internal/models"
	"route256/libs/logger"
)

func (r *CheckoutRepository) AddToCart(ctx context.Context, userID models.UserID, sku models.Sku, count models.Count) error {
	db := r.GetQueryEngine(ctx)
	childSpan, _ := opentracing.StartSpanFromContext(ctx, "repo: add_to_cart")
	defer childSpan.Finish()

	sq := squirrel.Expr(
		`INSERT INTO cart (user_id, item_sku, "count")
			VALUES($1, $2, $3)
			ON CONFLICT (user_id, item_sku) DO UPDATE SET count = $3;`,
		userID, sku, count)

	rawQuery, args, err := sq.ToSql()
	if err != nil {
		logger.Errorf("tosql error %v", err)
		return models.ErrInternal
	}
	var id []int
	if err = pgxscan.Select(ctx, db, &id, rawQuery, args...); err != nil {
		logger.Errorf("add to cart userID = %d, sku = %d, count = %d: database error: %v\n",
			userID, sku, count, err)
		return models.ErrInternal
	}

	return nil
}

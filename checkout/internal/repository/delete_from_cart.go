package repository

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/opentracing/opentracing-go"
	"route256/checkout/internal/models"
	"route256/libs/logger"
)

func (r *CheckoutRepository) DeleteFromCart(ctx context.Context, userID models.UserID, sku models.Sku, count models.Count) error {
	db := r.GetQueryEngine(ctx)
	childSpan, _ := opentracing.StartSpanFromContext(ctx, "repo: delete_from_cart")
	defer childSpan.Finish()
	sq := squirrel.Expr(
		`WITH deleted_row AS (
    	DELETE FROM cart
        WHERE item_sku = $1 AND user_id = $2 and count < $3
        RETURNING item_sku, user_id
		)
		UPDATE cart
		SET count = count - $3
		WHERE item_sku = $1 AND user_id = $2 AND NOT EXISTS (SELECT * FROM deleted_row);`,
		sku, userID, count)

	rawQuery, args, err := sq.ToSql()
	if err != nil {
		logger.Errorf("tosql error %v", err)
		return models.ErrInternal
	}
	var id []int
	if err = pgxscan.Select(ctx, db, &id, rawQuery, args...); err != nil {
		logger.Errorf("deleting from cart user=%d, sku=%d, count=%d: database error: %v\n", userID, sku, count, err)
		return models.ErrInternal
	}
	return nil
}

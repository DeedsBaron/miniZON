package repository

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/opentracing/opentracing-go"
	"route256/checkout/internal/models"
	"route256/checkout/internal/repository/schema"
	"route256/libs/logger"
)

func (r *CheckoutRepository) ListCart(ctx context.Context, userID models.UserID) (*models.Cart, error) {
	db := r.GetQueryEngine(ctx)
	childSpan, _ := opentracing.StartSpanFromContext(ctx, "repo: list_cart")
	defer childSpan.Finish()
	sq := squirrel.Expr(`
		select item_sku, "count"
		from cart
		where user_id = $1 and count > 0`, userID)

	rawQuery, args, err := sq.ToSql()
	if err != nil {
		logger.Errorf("tosql error %v", err)
		return nil, models.ErrInternal
	}
	var cartItems []schema.CartItem
	if err = pgxscan.Select(ctx, db, &cartItems, rawQuery, args...); err != nil {
		logger.Errorf("listing cart user=%d, database error: %v\n", userID, err)
		return nil, models.ErrInternal
	}
	if cartItems == nil {
		return nil, models.ErrNotFound
	}
	cart := models.Cart{
		Items:      make([]models.Item, 0, len(cartItems)),
		TotalPrice: 0,
	}
	for _, cartItem := range cartItems {
		cart.Items = append(cart.Items, models.Item{
			Sku:   models.Sku(cartItem.ItemSku),
			Count: models.Count(cartItem.Count),
			Name:  "",
			Price: 0,
		})
	}

	return &cart, nil
}

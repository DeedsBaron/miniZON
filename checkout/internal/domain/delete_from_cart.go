package domain

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"route256/checkout/internal/models"
	"route256/libs/logger"
)

func (b *Domain) DeleteFromCart(ctx context.Context, userID models.UserID, sku models.Sku, count models.Count) error {
	childSpan, childCtx := opentracing.StartSpanFromContext(ctx, "domain: delete_from_cart")
	defer childSpan.Finish()

	err := b.repo.DeleteFromCart(childCtx, userID, sku, count)
	if err != nil {
		logger.Errorw("creating order",
			"err", err.Error(),
			"component", "domain")
		return errors.WithMessage(errors.WithMessage(err, "deleteFromCart"), "domain")
	}
	return nil
}

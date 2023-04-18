package domain

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"route256/checkout/internal/models"
	"route256/libs/logger"
)

var (
	ErrLomsServiceStocks  = errors.New("loms service cant find stocks")
	ErrInsufficientStocks = errors.New("insufficient stocks")
)

func (b *Domain) AddToCart(ctx context.Context, userID models.UserID, sku models.Sku, count models.Count) error {
	childSpan, childCtx := opentracing.StartSpanFromContext(ctx, "domain: add_to_cart")
	defer childSpan.Finish()

	stocks, err := b.lomsService.Stocks(childCtx, uint32(sku))
	if err != nil {
		logger.Errorw("loms service",
			"err", err,
			"component", "domain")
		return ErrLomsServiceStocks
	}

	counter := int64(count)
	for _, stock := range stocks {
		counter -= int64(stock.Count)
		if counter <= 0 {
			err = b.tm.RunRepeteableReade(childCtx, func(ctxTX context.Context) error {
				var err error

				err = b.repo.AddToCart(ctxTX, userID, sku, count)
				if err != nil {
					return errors.WithMessage(err, "AddToCart")
				}
				return nil
			})

			if err != nil {
				logger.Errorw("adding to cart",
					"err", err.Error(),
					"component", "domain")
				return errors.WithMessage(err, "domain")
			}

			return nil
		}
	}
	return ErrInsufficientStocks
}

package domain

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"route256/loms/internal/models"

	"github.com/pkg/errors"
)

func (b *Domain) GetStocks(ctx context.Context, sku uint32) (*models.Stocks, error) {
	stocks := &models.Stocks{}

	childSpan, childCtx := opentracing.StartSpanFromContext(ctx, "domain: create_order")
	defer childSpan.Finish()

	err := b.tm.RunRepeteableReade(childCtx, func(ctxTX context.Context) error {
		var err error

		stocks, err = b.repo.GetStocks(ctxTX, sku)
		if err != nil {
			return errors.WithMessage(err, "getStocks")
		}
		return nil
	})

	if err != nil {
		return nil, errors.WithMessage(err, "domain")
	}
	return stocks, nil
}

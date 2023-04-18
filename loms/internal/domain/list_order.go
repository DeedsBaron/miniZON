package domain

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"route256/libs/logger"
	"route256/loms/internal/models"

	"github.com/pkg/errors"
)

func (b *Domain) ListOrder(ctx context.Context, orderID models.OrderID) (*models.Order, error) {
	childSpan, childCtx := opentracing.StartSpanFromContext(ctx, "domain: list_order")
	defer childSpan.Finish()
	orderList, err := b.repo.ListOrder(childCtx, orderID)
	if err != nil {
		logger.Infof("listing the order",
			"err", err.Error(),
			"component", "domain")
		return nil, errors.WithMessage(errors.WithMessage(err, "listOrder"), "domain")
	}
	return orderList, nil
}

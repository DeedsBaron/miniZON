package domain

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"route256/checkout/internal/models"
	"route256/libs/logger"
)

var (
	ErrLomsServiceCreateOrder = errors.New("loms service error")
)

func (b *Domain) Purchase(ctx context.Context, userID models.UserID) (models.OrderID, error) {
	childSpan, childCtx := opentracing.StartSpanFromContext(ctx, "domain: create_order")
	defer childSpan.Finish()
	cart, err := b.repo.ListCart(childCtx, userID)
	if err != nil {
		logger.Errorw("listing cart",
			"err", err.Error(),
			"component", "domain")
		return 0, errors.WithMessage(errors.WithMessage(err, "purchase"), "domain")
	}

	orderID, err := b.lomsService.CreateOrder(childCtx, userID, cart.Items)
	if err != nil {
		logger.Errorw("loms service",
			"err", err.Error(),
			"component", "domain")
		return 0, errors.WithMessage(errors.WithMessage(ErrLomsServiceCreateOrder, "purchase"), "domain")
	}

	return orderID, nil
}

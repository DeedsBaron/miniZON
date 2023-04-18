package domain

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"route256/libs/logger"
	"route256/loms/internal/models"

	"github.com/pkg/errors"
)

func (b *Domain) CancelOrder(ctx context.Context, orderID models.OrderID) error {
	childSpan, childCtx := opentracing.StartSpanFromContext(ctx, "domain: cancel_order")
	defer childSpan.Finish()

	err := b.tm.RunRepeteableReade(childCtx, func(ctxTX context.Context) error {
		var err error

		err = b.repo.CancelOrder(ctxTX, orderID)
		if err != nil {
			return errors.WithMessage(err, "cancelOrder")
		}
		err = b.repo.UpdateOrderStatus(ctxTX, []models.OrderID{orderID}, models.StatusCanceled)
		if err != nil {
			return errors.WithMessage(err, "updateOrderStatus")

		}
		oldStatus, err := b.repo.GetOrderStatus(ctx, orderID)
		if err != nil {
			return errors.WithMessage(err, "getOrderStatus")

		}
		err = b.repo.SetOutbox(ctxTX, orderID, oldStatus, models.StatusCanceled)
		if err != nil {
			return errors.WithMessage(err, "setOutbox")
		}

		return nil
	})

	if err != nil {
		logger.Errorw("canceling order",
			"err", err.Error(),
			"component", "domain")
		return errors.WithMessage(err, "domain")
	}
	return nil
}

package domain

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"route256/libs/logger"
	"route256/loms/internal/models"

	"github.com/pkg/errors"
)

func (b *Domain) OrderPayed(ctx context.Context, orderID models.OrderID) error {
	childSpan, childCtx := opentracing.StartSpanFromContext(ctx, "domain: pay_order")
	defer childSpan.Finish()
	err := b.tm.RunRepeteableReade(childCtx, func(ctxTX context.Context) error {
		var err error

		err = b.repo.ReduceStock(ctxTX, orderID)
		if err != nil {
			return errors.WithMessage(err, "reduceStock")
		}

		err = b.repo.UpdateOrderStatus(ctxTX, []models.OrderID{orderID}, models.StatusPayed)
		if err != nil {
			return errors.WithMessage(err, "updateOrderStatus")
		}
		oldStatus, err := b.repo.GetOrderStatus(ctx, orderID)
		if err != nil {
			return errors.WithMessage(err, "getOrderStatus")

		}
		err = b.repo.SetOutbox(ctxTX, orderID, oldStatus, models.StatusPayed)
		if err != nil {
			return errors.WithMessage(err, "setOutbox")
		}

		return nil
	})

	if err != nil {
		logger.Infof("paying the order",
			"err", err.Error(),
			"component", "domain")
		return errors.WithMessage(err, "domain")
	}
	return nil
}

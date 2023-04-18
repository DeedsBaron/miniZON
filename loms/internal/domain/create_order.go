package domain

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"route256/libs/logger"
	"route256/loms/internal/models"

	"github.com/pkg/errors"
)

func (b *Domain) CreateOrder(ctx context.Context, createOrder models.CreateOrder) (models.OrderID, error) {
	var createdOrderID models.OrderID

	childSpan, childCtx := opentracing.StartSpanFromContext(ctx, "domain: create_order")
	defer childSpan.Finish()

	err := b.tm.RunRepeteableReade(childCtx, func(ctxTX context.Context) error {
		var err error

		createdOrderID, err = b.repo.CreateOrder(ctxTX, createOrder)
		if err != nil {
			return errors.WithMessage(err, "createOrder")
		}

		err = b.repo.SetOutbox(ctxTX, createdOrderID, models.StatusNone, models.StatusNew)
		if err != nil {
			return errors.WithMessage(err, "setOutbox")
		}

		for _, item := range createOrder.Items {
			err = b.ReserveStocks(ctxTX, createdOrderID, item.Sku, item.Count)
			if err != nil {
				updateErr := b.repo.UpdateOrderStatus(ctxTX, []models.OrderID{createdOrderID}, models.StatusFailed)
				if updateErr != nil {
					return errors.WithMessage(updateErr, "updateOrderStatus after failed reservation")
				}

				setErr := b.repo.SetOutbox(ctxTX, createdOrderID, models.StatusNew, models.StatusFailed)
				if setErr != nil {
					return errors.WithMessage(err, "setOutbox after failed reservation")
				}

				return errors.WithMessage(err, "reserveStocks")
			}
		}
		err = b.repo.UpdateOrderStatus(ctxTX, []models.OrderID{createdOrderID}, models.StatusAwaitingPayment)
		if err != nil {
			return errors.WithMessage(err, "updateOrderStatus")
			return err
		}

		err = b.repo.SetOutbox(ctxTX, createdOrderID, models.StatusNew, models.StatusAwaitingPayment)
		if err != nil {
			return errors.WithMessage(err, "setOutbox")
		}

		return nil
	})

	if err != nil {
		logger.Errorw("creating order",
			"err", err.Error(),
			"component", "domain")
		return 0, errors.WithMessage(err, "domain")
	}
	return createdOrderID, nil
}

package repository

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"route256/loms/internal/models"
)

func (r *LomsRepository) CancelOrder(ctx context.Context, orderID models.OrderID) error {
	status, err := r.GetOrderStatus(ctx, orderID)
	childSpan, _ := opentracing.StartSpanFromContext(ctx, "repo: cancel_order")
	defer childSpan.Finish()
	if err != nil {
		return err
	}

	switch status {
	case models.StatusPayed:
		err = r.IncreaseStock(ctx, orderID)
		if err != nil {
			return err
		}
		fallthrough
	default:
		err = r.UpdateOrderStatus(ctx, []models.OrderID{orderID}, models.StatusCanceled)
		if err != nil {
			return err
		}
	}
	return nil
}

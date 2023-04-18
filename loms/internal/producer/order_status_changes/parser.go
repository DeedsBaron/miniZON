package order_status_changes

import (
	"fmt"
	"route256/loms/internal/models"
	"route256/loms/internal/pb/kafka/orders_status_changes"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type OrdersStatusChangesMessage struct {
	OldStatus string
	NewStatus string
	ChangedAt time.Time
}

func NewMessage(outbox models.Outbox) (string, []byte, error) {
	msg := &orders_status_changes.OrderStatusChanges{
		OldStatus: string(outbox.OldStatus),
		NewStatus: string(outbox.NewStatus),
		ChangedAt: timestamppb.New(outbox.ChangedAt),
	}

	value, err := protojson.Marshal(msg)
	if err != nil {
		return "", nil, errors.WithMessage(err, "proto marshalling")
	}
	key := fmt.Sprint(outbox.OrderID)

	return key, value, nil
}

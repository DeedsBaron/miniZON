package models

import "time"

type OutboxID int

type Outbox struct {
	ID        OutboxID
	OrderID   OrderID
	OldStatus Status
	NewStatus Status
	ChangedAt time.Time
	IsSent    bool
}

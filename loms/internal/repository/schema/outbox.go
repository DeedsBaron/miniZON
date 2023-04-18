package schema

import "time"

type Outbox struct {
	ID        int       `db:"id"`
	OrderID   int       `db:"order_id"`
	OldStatus string    `db:"old_status"`
	NewStatus string    `db:"new_status"`
	ChangedAt time.Time `db:"changed_at"`
}

package models

type CreateOrder struct {
	User  int64
	Items []OrderItem
}

type OrderItem struct {
	Sku   uint32
	Count uint32
}

type OrderID int64

type Status string

const (
	StatusNew             Status = "new"
	StatusAwaitingPayment Status = "awaiting_payment"
	StatusFailed          Status = "failed"
	StatusPayed           Status = "payed"
	StatusCanceled        Status = "cancelled"
	StatusNone            Status = "none"
)

type Order struct {
	Status Status
	User   User
	Items  []OrderItem
}

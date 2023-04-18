package repository

import (
	"context"

	"route256/libs/transactor"
	"route256/loms/internal/models"
)

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate minimock -i Repository -o ./mocks/ -s "_minimock.go"

type Repository interface {
	GetStocks(ctx context.Context, sku uint32) (*models.Stocks, error)
	CreateOrder(ctx context.Context, createOrder models.CreateOrder) (models.OrderID, error)
	ReserveStocks(ctx context.Context, orderId models.OrderID, itemSku uint32, needToReserveCount uint64, warehouseId models.WarehouseID) error
	UpdateOrderStatus(ctx context.Context, orderID []models.OrderID, status models.Status) error
	ListOrder(ctx context.Context, orderID models.OrderID) (*models.Order, error)
	GetOrderStatus(ctx context.Context, orderID models.OrderID) (models.Status, error)
	GetOrderUser(ctx context.Context, orderID models.OrderID) (models.User, error)
	ReduceStock(ctx context.Context, orderID models.OrderID) error
	CancelOrder(ctx context.Context, orderID models.OrderID) error
	GetUnpayedOrdersWithinTimeout(ctx context.Context) ([]models.OrderID, error)
	SetOutbox(ctx context.Context, orderID models.OrderID, oldStatus, newStatus models.Status) error
	GetOutbox(ctx context.Context) ([]models.Outbox, error)
	UpdateOutbox(ctx context.Context, ids []models.OutboxID) error
}

type LomsRepository struct {
	transactor.QueryEngineProvider
}

func NewLomsRepository(provider transactor.QueryEngineProvider) *LomsRepository {
	return &LomsRepository{
		provider,
	}
}

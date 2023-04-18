package domain

import (
	"context"
	"route256/libs/transactor"
	"route256/loms/internal/models"
	"route256/loms/internal/repository"
)

type BusinessLogic interface {
	GetStocks(ctx context.Context, sku uint32) (*models.Stocks, error)
	CreateOrder(ctx context.Context, createOrder models.CreateOrder) (models.OrderID, error)
	ListOrder(ctx context.Context, orderID models.OrderID) (*models.Order, error)
	OrderPayed(ctx context.Context, orderID models.OrderID) error
	CancelOrder(ctx context.Context, orderID models.OrderID) error
}

type Domain struct {
	repo repository.Repository
	tm   transactor.TransactionManager
}

func NewBuisnessLogic(repo repository.Repository, tm transactor.TransactionManager) *Domain {
	return &Domain{
		repo: repo,
		tm:   tm,
	}
}

func NewMockBuisnessLogic(deps ...interface{}) *Domain {
	d := Domain{}

	for _, v := range deps {
		switch s := v.(type) {
		case repository.Repository:
			d.repo = s
		case transactor.TransactionManager:
			d.tm = s
		}
	}

	return &d
}

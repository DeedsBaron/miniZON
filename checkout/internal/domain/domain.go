package domain

import (
	"context"

	"route256/checkout/internal/clients/loms"
	"route256/checkout/internal/clients/ps"
	"route256/checkout/internal/models"
	"route256/checkout/internal/repository"
	"route256/libs/cache"
	"route256/libs/transactor"
)

type BusinessLogic interface {
	AddToCart(ctx context.Context, userID models.UserID, sku models.Sku, count models.Count) error
	DeleteFromCart(ctx context.Context, userID models.UserID, sku models.Sku, count models.Count) error
	ListCart(ctx context.Context, userID models.UserID) (*models.Cart, error)
	Purchase(ctx context.Context, userID models.UserID) (models.OrderID, error)
	CacheDelete(ctx context.Context, keys models.Keys) error
}

type Domain struct {
	repo           repository.Repository
	lomsService    loms.LomsService
	productService ps.ProductService
	tm             transactor.TransactionManager
	cache          cache.Cache[string]
}

func NewBuisnessLogic(lomsService loms.LomsService,
	productService ps.ProductService,
	repo repository.Repository,
	tm transactor.TransactionManager,
	cache cache.Cache[string]) *Domain {
	return &Domain{
		lomsService:    lomsService,
		productService: productService,
		repo:           repo,
		tm:             tm,
		cache:          cache,
	}
}

func NewMockBuisnessLogic(deps ...interface{}) *Domain {
	d := Domain{}

	for _, v := range deps {
		switch s := v.(type) {
		case ps.ProductService:
			d.productService = s
		case loms.LomsService:
			d.lomsService = s
		case repository.Repository:
			d.repo = s
		case transactor.TransactionManager:
			d.tm = s
		}
	}

	return &d
}

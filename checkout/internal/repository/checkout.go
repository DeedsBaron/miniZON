package repository

import (
	"context"

	"route256/checkout/internal/models"
	"route256/libs/transactor"
)

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate minimock -i Repository -o ./mocks/ -s "_minimock.go"

type Repository interface {
	AddToCart(ctx context.Context, userID models.UserID, sku models.Sku, count models.Count) error
	DeleteFromCart(ctx context.Context, userID models.UserID, sku models.Sku, count models.Count) error
	ListCart(ctx context.Context, userID models.UserID) (*models.Cart, error)
}

type CheckoutRepository struct {
	transactor.QueryEngineProvider
}

func NewCheckoutRepository(provider transactor.QueryEngineProvider) *CheckoutRepository {
	return &CheckoutRepository{
		provider,
	}
}

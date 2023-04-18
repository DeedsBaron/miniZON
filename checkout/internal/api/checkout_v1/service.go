package checkout_v1

import (
	"route256/checkout/internal/domain"
	desc "route256/checkout/pkg/checkout_v1"
)

type Implementation struct {
	desc.UnimplementedCheckoutV1Server
	domain.BusinessLogic
}

func NewCheckoutV1(implementation domain.BusinessLogic) *Implementation {
	return &Implementation{
		desc.UnimplementedCheckoutV1Server{},
		implementation,
	}
}

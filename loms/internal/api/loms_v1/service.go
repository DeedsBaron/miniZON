package loms_v1

import (
	"route256/loms/internal/domain"
	desc "route256/loms/pkg/loms_v1"
)

type Implementation struct {
	desc.UnimplementedLomsV1Server
	domain.BusinessLogic
}

func NewLomsV1(implementation domain.BusinessLogic) *Implementation {
	return &Implementation{
		desc.UnimplementedLomsV1Server{},
		implementation,
	}
}

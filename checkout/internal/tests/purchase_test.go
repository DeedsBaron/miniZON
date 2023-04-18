package tests

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"route256/checkout/internal/clients/loms"
	lomsMock "route256/checkout/internal/clients/loms/mocks"
	domain2 "route256/checkout/internal/domain"
	"route256/checkout/internal/models"
	"route256/checkout/internal/repository"
	repoMock "route256/checkout/internal/repository/mocks"
)

func TestPurchase(t *testing.T) {
	t.Parallel()
	type checkoutRepositoryMockFunc func(mc *minimock.Controller) repository.Repository
	type lomsClientMockFunc func(mc *minimock.Controller) loms.LomsService

	type args struct {
		ctx    context.Context
		userID models.UserID
	}

	var (
		mc      = minimock.NewController(t)
		ctx     = context.Background()
		repoErr = errors.New("repo error")
		lomsErr = errors.New("loms service error")
		user    = models.UserID(gofakeit.Int64())
		sku     = gofakeit.Uint32()
		count   = gofakeit.Uint32()
		orderID = models.OrderID(gofakeit.Int64())

		checkoutRepoReq = user
		checkoutRepoRes = &models.Cart{
			Items: []models.Item{
				models.Item{
					Sku:   models.Sku(sku),
					Count: models.Count(count),
					Name:  "",
					Price: 0,
				},
			},
			TotalPrice: 0,
		}

		lomsClientResp = &orderID

		expectedRes = &orderID
	)
	t.Cleanup(mc.Finish)

	tests := []struct {
		name                   string
		args                   args
		want                   *models.OrderID
		err                    error
		checkoutRepositoryMock checkoutRepositoryMockFunc
		lomsClientMock         lomsClientMockFunc
	}{
		{
			name: "positive case",
			args: args{
				ctx:    ctx,
				userID: user,
			},
			want: expectedRes,
			err:  nil,
			checkoutRepositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := repoMock.NewRepositoryMock(mc)
				mock.ListCartMock.Expect(ctx, checkoutRepoReq).Return(checkoutRepoRes, nil)
				return mock
			},
			lomsClientMock: func(mc *minimock.Controller) loms.LomsService {
				mock := lomsMock.NewLomsServiceMock(mc)
				mock.CreateOrderMock.Expect(ctx, user, checkoutRepoRes.Items).Return(lomsClientResp, nil)
				return mock
			},
		},
		{
			name: "repo error",
			args: args{
				ctx:    ctx,
				userID: user,
			},
			want: nil,
			err:  repoErr,
			checkoutRepositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := repoMock.NewRepositoryMock(mc)
				mock.ListCartMock.Expect(ctx, checkoutRepoReq).Return(nil, repoErr)
				return mock
			},
			lomsClientMock: func(mc *minimock.Controller) loms.LomsService {
				mock := lomsMock.NewLomsServiceMock(mc)
				return mock
			},
		},
		{
			name: "loms client error",
			args: args{
				ctx:    ctx,
				userID: user,
			},
			want: nil,
			err:  lomsErr,
			checkoutRepositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := repoMock.NewRepositoryMock(mc)
				mock.ListCartMock.Expect(ctx, checkoutRepoReq).Return(checkoutRepoRes, nil)
				return mock
			},
			lomsClientMock: func(mc *minimock.Controller) loms.LomsService {
				mock := lomsMock.NewLomsServiceMock(mc)
				mock.CreateOrderMock.Expect(ctx, user, checkoutRepoRes.Items).Return(nil, lomsErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			domain := domain2.NewMockBuisnessLogic(tt.checkoutRepositoryMock(mc), tt.lomsClientMock(mc))
			res, err := domain.Purchase(tt.args.ctx, tt.args.userID)
			require.Equal(t, tt.want, res)
			if tt.err != nil {
				require.ErrorContains(t, err, tt.err.Error())
			} else {
				require.Equal(t, tt.err, err)
			}
		})
	}
}

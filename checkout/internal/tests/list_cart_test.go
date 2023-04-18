package tests

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"route256/checkout/internal/clients/ps"
	psMock "route256/checkout/internal/clients/ps/mocks"
	"route256/checkout/internal/config"
	domain2 "route256/checkout/internal/domain"
	"route256/checkout/internal/models"
	"route256/checkout/internal/repository"
	repoMock "route256/checkout/internal/repository/mocks"
)

func TestListCart(t *testing.T) {
	t.Parallel()
	type checkoutRepositoryMockFunc func(mc *minimock.Controller) repository.Repository
	type psClientMockFunc func(mc *minimock.Controller) ps.ProductService

	type args struct {
		ctx    context.Context
		userID models.UserID
	}

	var (
		mc      = minimock.NewController(t)
		ctx     = context.Background()
		repoErr = errors.New("repo error")
		psErr   = errors.New("product service error")
		user    = models.UserID(gofakeit.Int64())
		sku     = gofakeit.Uint32()
		count   = gofakeit.Uint32()
		name    = gofakeit.Name()
		price   = gofakeit.Uint32()

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
		psClientReq  = sku
		psClientResp = &models.ProductInfo{
			Name:      name,
			Price:     price,
			CartIndex: 0,
		}

		expectedRes = models.Cart{
			Items: []models.Item{
				models.Item{
					Sku:   models.Sku(sku),
					Count: models.Count(count),
					Name:  name,
					Price: price,
				},
			},
			TotalPrice: price,
		}
	)
	t.Cleanup(mc.Finish)

	tests := []struct {
		name                   string
		args                   args
		want                   *models.Cart
		err                    error
		checkoutRepositoryMock checkoutRepositoryMockFunc
		psClientMock           psClientMockFunc
	}{
		{
			name: "positive case",
			args: args{
				ctx:    ctx,
				userID: user,
			},
			want: &expectedRes,
			err:  nil,
			checkoutRepositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := repoMock.NewRepositoryMock(mc)
				mock.ListCartMock.Expect(ctx, checkoutRepoReq).Return(checkoutRepoRes, nil)
				return mock
			},
			psClientMock: func(mc *minimock.Controller) ps.ProductService {
				mock := psMock.NewProductServiceMock(mc)
				mock.GetProductMock.Expect(ctx, psClientReq).Return(psClientResp, nil)
				return mock
			},
		},
		{
			name: "list cart repo error",
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
			psClientMock: func(mc *minimock.Controller) ps.ProductService {
				mock := psMock.NewProductServiceMock(mc)
				return mock
			},
		},
		{
			name: "product service client error",
			args: args{
				ctx:    ctx,
				userID: user,
			},
			want: nil,
			err:  psErr,
			checkoutRepositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := repoMock.NewRepositoryMock(mc)
				mock.ListCartMock.Expect(ctx, checkoutRepoReq).Return(checkoutRepoRes, nil)
				return mock
			},
			psClientMock: func(mc *minimock.Controller) ps.ProductService {
				mock := psMock.NewProductServiceMock(mc)
				mock.GetProductMock.Expect(ctx, psClientReq).Return(nil, psErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			config.Data.Workers.Amount = 5
			domain := domain2.NewMockBuisnessLogic(tt.checkoutRepositoryMock(mc), tt.psClientMock(mc))
			res, err := domain.ListCart(tt.args.ctx, tt.args.userID)
			require.Equal(t, tt.want, res)
			if tt.err != nil {
				require.ErrorContains(t, err, tt.err.Error())
			} else {
				require.Equal(t, tt.err, err)
			}
		})
	}
}

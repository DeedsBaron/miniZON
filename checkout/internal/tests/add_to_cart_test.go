package tests

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"route256/checkout/internal/clients/loms"
	lomsMock "route256/checkout/internal/clients/loms/mocks"
	domain2 "route256/checkout/internal/domain"
	"route256/checkout/internal/models"
	"route256/checkout/internal/repository"
	repoMock "route256/checkout/internal/repository/mocks"
	"route256/libs/transactor"
	pgxMocks "route256/libs/transactor/mocks"
)

func TestAddToCart(t *testing.T) {
	t.Parallel()
	type lomsClientMockFunc func(mc *minimock.Controller) loms.LomsService
	type dbClientMockFunc func(mc *minimock.Controller) transactor.DbClient
	type checkoutRepositoryMockFunc func(mc *minimock.Controller) repository.Repository

	type args struct {
		ctx    context.Context
		userID models.UserID
		sku    models.Sku
		count  models.Count
	}

	var (
		tx     = pgxMocks.NewTxMock(t)
		mc     = minimock.NewController(t)
		ctx    = context.Background()
		ctxTx  = context.WithValue(ctx, transactor.Key, tx)
		txOpts = pgx.TxOptions{
			IsoLevel: pgx.RepeatableRead,
		}
		lomsErr     = errors.New("loms service cant find stocks")
		user        = models.UserID(gofakeit.Int64())
		sku         = models.Sku(gofakeit.Uint32())
		warehouseID = gofakeit.Int64()
		count       = models.Count(100)

		lomsClientReq  = sku
		lomsClientResp = []models.Stock{
			models.Stock{
				WarehouseID: warehouseID,
				Count:       uint64(count),
			},
		}
	)
	t.Cleanup(mc.Finish)

	tests := []struct {
		name                   string
		args                   args
		want                   error
		err                    error
		lomsClientMock         lomsClientMockFunc
		dbClientMock           dbClientMockFunc
		checkoutRepositoryMock checkoutRepositoryMockFunc
	}{
		{
			name: "positive case",
			args: args{
				ctx:    ctx,
				userID: user,
				sku:    sku,
				count:  count,
			},
			want: nil,
			err:  nil,
			lomsClientMock: func(mc *minimock.Controller) loms.LomsService {
				mock := lomsMock.NewLomsServiceMock(mc)
				mock.StocksMock.Expect(ctx, uint32(lomsClientReq)).Return(lomsClientResp, nil)
				return mock
			},
			dbClientMock: func(mc *minimock.Controller) transactor.DbClient {
				mock := pgxMocks.NewDbClientMock(mc)
				mock.BeginTxMock.Expect(ctx, txOpts).Return(tx, nil)
				tx.CommitMock.Expect(ctx).Return(nil)
				return mock
			},
			checkoutRepositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := repoMock.NewRepositoryMock(mc)
				mock.AddToCartMock.Expect(ctxTx, user, sku, count).Return(nil)
				return mock
			},
		},
		{
			name: "loms service stock error",
			args: args{
				ctx:    ctx,
				userID: user,
				sku:    sku,
				count:  count,
			},
			want: lomsErr,
			err:  lomsErr,
			lomsClientMock: func(mc *minimock.Controller) loms.LomsService {
				mock := lomsMock.NewLomsServiceMock(mc)
				mock.StocksMock.Expect(ctx, uint32(lomsClientReq)).Return(nil, lomsErr)
				return mock
			},
			dbClientMock: func(mc *minimock.Controller) transactor.DbClient {
				mock := pgxMocks.NewDbClientMock(mc)
				return mock
			},
			checkoutRepositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := repoMock.NewRepositoryMock(mc)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tm := transactor.NewTransactor(tt.dbClientMock(mc))

			domain := domain2.NewMockBuisnessLogic(tt.lomsClientMock(mc), tt.checkoutRepositoryMock(mc), tm)
			err := domain.AddToCart(tt.args.ctx, tt.args.userID, tt.args.sku, tt.args.count)
			if tt.err != nil {
				require.ErrorContains(t, err, tt.err.Error())
			} else {
				require.Equal(t, tt.err, err)
			}
		})
	}
}

package tests

import (
	"context"
	"route256/libs/transactor"
	pgxMocks "route256/libs/transactor/mocks"
	domain2 "route256/loms/internal/domain"
	"route256/loms/internal/models"
	"route256/loms/internal/repository"
	repoMock "route256/loms/internal/repository/mocks"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestCreateOrder(t *testing.T) {
	t.Parallel()
	type lomsRepositoryMockFunc func(mc *minimock.Controller) repository.Repository
	type transactorManagerMockFunc func(mc *minimock.Controller) transactor.TransactionManager
	type dbClientMockFunc func(mc *minimock.Controller) transactor.DbClient

	type args struct {
		ctx         context.Context
		createOrder models.CreateOrder
	}

	var (
		tx      = pgxMocks.NewTxMock(t)
		mc      = minimock.NewController(t)
		ctx     = context.Background()
		ctxTx   = context.WithValue(ctx, transactor.Key, tx)
		repoErr = errors.New("repo error")
		orderID = models.OrderID(gofakeit.Int64())
		user    = models.User(gofakeit.Int64())
		items   = []models.OrderItem{models.OrderItem{
			Sku:   gofakeit.Uint32(),
			Count: uint32(15),
		}}
		warehouseID = models.WarehouseID(gofakeit.Uint32())
		count       = uint64(100)
		stocks      = models.Stocks{
			Stocks: []models.StocksItem{
				{
					WarehouseID: models.WarehouseID(warehouseID),
					Count:       count,
				},
			},
		}
		createOrder = models.CreateOrder{
			User:  int64(user),
			Items: items,
		}
		txOpts = pgx.TxOptions{
			IsoLevel: pgx.RepeatableRead,
		}
		lomsRepoReq  = createOrder
		lomsRepoResp = orderID
		expectedRes  = orderID
	)
	t.Cleanup(mc.Finish)

	tests := []struct {
		name               string
		args               args
		want               *models.OrderID
		err                error
		lomsRepositoryMock lomsRepositoryMockFunc
		transactorMock     transactorManagerMockFunc
		dbClientMock       dbClientMockFunc
	}{
		{
			name: "positive case",
			args: args{
				ctx:         ctx,
				createOrder: createOrder,
			},
			want: &expectedRes,
			err:  nil,
			lomsRepositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := repoMock.NewRepositoryMock(mc)
				mock.GetStocksMock.Expect(ctxTx, items[0].Sku).Return(&stocks, nil)
				mock.CreateOrderMock.Expect(ctxTx, lomsRepoReq).Return(lomsRepoResp, nil)
				mock.ReserveStocksMock.Expect(ctxTx, orderID, items[0].Sku, 15, warehouseID).Return(nil)
				mock.UpdateOrderStatusMock.Expect(ctxTx, []models.OrderID{orderID}, models.StatusAwaitingPayment).Return(nil)
				return mock
			},
			dbClientMock: func(mc *minimock.Controller) transactor.DbClient {
				mock := pgxMocks.NewDbClientMock(mc)
				mock.BeginTxMock.Expect(ctx, txOpts).Return(tx, nil)
				tx.CommitMock.Expect(ctx).Return(nil)
				return mock
			},
		},
		{
			name: "repo create order error",
			args: args{
				ctx:         ctx,
				createOrder: createOrder,
			},
			want: nil,
			err:  repoErr,
			lomsRepositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := repoMock.NewRepositoryMock(mc)
				mock.CreateOrderMock.Expect(ctxTx, lomsRepoReq).Return(0, repoErr)
				return mock
			},
			dbClientMock: func(mc *minimock.Controller) transactor.DbClient {
				mock := pgxMocks.NewDbClientMock(mc)
				mock.BeginTxMock.Expect(ctx, txOpts).Return(tx, nil)
				tx.CommitMock.Expect(ctx).Return(nil)
				tx.RollbackMock.Expect(ctx).Return(repoErr)
				return mock
			},
		},
		{
			name: "reserve stocks get stocks error",
			args: args{
				ctx:         ctx,
				createOrder: createOrder,
			},
			want: nil,
			err:  repoErr,
			lomsRepositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := repoMock.NewRepositoryMock(mc)
				mock.CreateOrderMock.Expect(ctxTx, lomsRepoReq).Return(lomsRepoResp, nil)
				mock.GetStocksMock.Expect(ctxTx, items[0].Sku).Return(nil, repoErr)
				mock.UpdateOrderStatusMock.Expect(ctxTx, []models.OrderID{orderID}, models.StatusFailed).Return(nil)
				return mock
			},
			dbClientMock: func(mc *minimock.Controller) transactor.DbClient {
				mock := pgxMocks.NewDbClientMock(mc)
				mock.BeginTxMock.Expect(ctx, txOpts).Return(tx, nil)
				tx.CommitMock.Expect(ctx).Return(nil)
				tx.RollbackMock.Expect(ctx).Return(repoErr)
				return mock
			},
		},
		{
			name: "reserve stocks repo error",
			args: args{
				ctx:         ctx,
				createOrder: createOrder,
			},
			want: nil,
			err:  repoErr,
			lomsRepositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := repoMock.NewRepositoryMock(mc)
				mock.GetStocksMock.Expect(ctxTx, items[0].Sku).Return(&stocks, nil)
				mock.CreateOrderMock.Expect(ctxTx, lomsRepoReq).Return(lomsRepoResp, nil)
				mock.ReserveStocksMock.Expect(ctxTx, orderID, items[0].Sku, 15, warehouseID).Return(repoErr)
				mock.UpdateOrderStatusMock.Expect(ctxTx, []models.OrderID{orderID}, models.StatusFailed).Return(nil)
				return mock
			},
			dbClientMock: func(mc *minimock.Controller) transactor.DbClient {
				mock := pgxMocks.NewDbClientMock(mc)
				mock.BeginTxMock.Expect(ctx, txOpts).Return(tx, nil)
				tx.CommitMock.Expect(ctx).Return(nil)
				tx.RollbackMock.Expect(ctx).Return(repoErr)
				return mock
			},
		},
		{
			name: "update status after reservation error",
			args: args{
				ctx:         ctx,
				createOrder: createOrder,
			},
			want: nil,
			err:  repoErr,
			lomsRepositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := repoMock.NewRepositoryMock(mc)
				mock.CreateOrderMock.Expect(ctxTx, lomsRepoReq).Return(lomsRepoResp, nil)
				mock.GetStocksMock.Expect(ctxTx, items[0].Sku).Return(nil, repoErr)
				mock.UpdateOrderStatusMock.Expect(ctxTx, []models.OrderID{orderID}, models.StatusFailed).Return(repoErr)
				return mock
			},
			dbClientMock: func(mc *minimock.Controller) transactor.DbClient {
				mock := pgxMocks.NewDbClientMock(mc)
				mock.BeginTxMock.Expect(ctx, txOpts).Return(tx, nil)
				tx.CommitMock.Expect(ctx).Return(nil)
				tx.RollbackMock.Expect(ctx).Return(repoErr)
				return mock
			},
		},
		{
			name: "update order status error",
			args: args{
				ctx:         ctx,
				createOrder: createOrder,
			},
			want: nil,
			err:  repoErr,
			lomsRepositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := repoMock.NewRepositoryMock(mc)
				mock.GetStocksMock.Expect(ctxTx, items[0].Sku).Return(&stocks, nil)
				mock.CreateOrderMock.Expect(ctxTx, lomsRepoReq).Return(lomsRepoResp, nil)
				mock.ReserveStocksMock.Expect(ctxTx, orderID, items[0].Sku, 15, warehouseID).Return(nil)
				mock.UpdateOrderStatusMock.Expect(ctxTx, []models.OrderID{orderID}, models.StatusAwaitingPayment).Return(repoErr)
				return mock
			},
			dbClientMock: func(mc *minimock.Controller) transactor.DbClient {
				mock := pgxMocks.NewDbClientMock(mc)
				mock.BeginTxMock.Expect(ctx, txOpts).Return(tx, nil)
				tx.CommitMock.Expect(ctx).Return(nil)
				tx.RollbackMock.Expect(ctx).Return(repoErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tm := transactor.NewTransactor(tt.dbClientMock(mc))

			domain := domain2.NewMockBuisnessLogic(tt.lomsRepositoryMock(mc), tm)
			res, err := domain.CreateOrder(tt.args.ctx, tt.args.createOrder)
			require.Equal(t, tt.want, res)
			if tt.err != nil {
				require.ErrorContains(t, err, tt.err.Error())
			} else {
				require.Equal(t, tt.err, err)
			}
		})
	}
}

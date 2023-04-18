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

func TestOrderPayed(t *testing.T) {
	t.Parallel()
	type lomsRepositoryMockFunc func(mc *minimock.Controller) repository.Repository
	type dbClientMockFunc func(mc *minimock.Controller) transactor.DbClient

	type args struct {
		ctx     context.Context
		orderID models.OrderID
	}

	var (
		tx      = pgxMocks.NewTxMock(t)
		mc      = minimock.NewController(t)
		ctx     = context.Background()
		ctxTx   = context.WithValue(ctx, transactor.Key, tx)
		repoErr = errors.New("repo error")
		txOpts  = pgx.TxOptions{
			IsoLevel: pgx.RepeatableRead,
		}
		orderID     = models.OrderID(gofakeit.Int64())
		lomsRepoReq = orderID
	)
	t.Cleanup(mc.Finish)

	tests := []struct {
		name               string
		args               args
		want               error
		err                error
		lomsRepositoryMock lomsRepositoryMockFunc
		dbClientMock       dbClientMockFunc
	}{
		{
			name: "positive case",
			args: args{
				ctx:     ctx,
				orderID: orderID,
			},
			want: nil,
			err:  nil,
			lomsRepositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := repoMock.NewRepositoryMock(mc)
				mock.ReduceStockMock.Expect(ctxTx, orderID).Return(nil)
				mock.UpdateOrderStatusMock.Expect(ctxTx, []models.OrderID{orderID}, models.StatusPayed).Return(nil)
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
			name: "repo reduce stock error",
			args: args{
				ctx:     ctx,
				orderID: orderID,
			},
			want: nil,
			err:  repoErr,
			lomsRepositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := repoMock.NewRepositoryMock(mc)
				mock.ReduceStockMock.Expect(ctxTx, lomsRepoReq).Return(repoErr)
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
			name: "repo update status error",
			args: args{
				ctx:     ctx,
				orderID: orderID,
			},
			want: nil,
			err:  repoErr,
			lomsRepositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := repoMock.NewRepositoryMock(mc)
				mock.ReduceStockMock.Expect(ctxTx, lomsRepoReq).Return(nil)
				mock.UpdateOrderStatusMock.Expect(ctxTx, []models.OrderID{orderID}, models.StatusPayed).Return(repoErr)
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
			err := domain.OrderPayed(tt.args.ctx, tt.args.orderID)
			if tt.err != nil {
				require.ErrorContains(t, err, tt.err.Error())
			} else {
				require.Equal(t, tt.err, err)
			}
		})
	}
}

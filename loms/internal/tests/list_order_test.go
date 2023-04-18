package tests

import (
	"context"
	domain2 "route256/loms/internal/domain"
	"route256/loms/internal/models"
	"route256/loms/internal/repository"
	"route256/loms/internal/repository/mocks"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestGetList(t *testing.T) {
	t.Parallel()
	type lomsRepositoryMockFunc func(mc *minimock.Controller) repository.Repository

	type args struct {
		ctx     context.Context
		orderID models.OrderID
	}

	var (
		mc      = minimock.NewController(t)
		ctx     = context.Background()
		repoErr = errors.New("repo error")
		status  = models.Status(gofakeit.RandomString([]string{
			string(models.StatusCanceled),
			string(models.StatusPayed),
			string(models.StatusFailed),
			string(models.StatusAwaitingPayment),
			string(models.StatusNew),
		}))
		user  = models.User(gofakeit.Int64())
		items = []models.OrderItem{
			{
				Sku:   gofakeit.Uint32(),
				Count: gofakeit.Uint32(),
			},
		}
		orderID = models.OrderID(gofakeit.Int64())

		lomsRepoReq = orderID

		repoRes = &models.Order{
			Status: status,
			User:   user,
			Items:  items,
		}

		expectedRes = &models.Order{
			Status: status,
			User:   user,
			Items:  items,
		}
	)
	t.Cleanup(mc.Finish)

	tests := []struct {
		name               string
		args               args
		want               *models.Order
		err                error
		lomsRepositoryMock lomsRepositoryMockFunc
	}{
		{
			name: "positive case",
			args: args{
				ctx:     ctx,
				orderID: orderID,
			},
			want: expectedRes,
			err:  nil,
			lomsRepositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := mocks.NewRepositoryMock(mc)
				mock.ListOrderMock.Expect(ctx, lomsRepoReq).Return(repoRes, nil)
				return mock
			},
		},
		{
			name: "negative case - repository error",
			args: args{
				ctx:     ctx,
				orderID: orderID,
			},
			want: nil,
			err:  repoErr,
			lomsRepositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := mocks.NewRepositoryMock(mc)
				mock.ListOrderMock.Expect(ctx, lomsRepoReq).Return(nil, repoErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			domain := domain2.NewMockBuisnessLogic(tt.lomsRepositoryMock(mc))
			res, err := domain.ListOrder(tt.args.ctx, tt.args.orderID)
			require.Equal(t, tt.want, res)
			if tt.err != nil {
				require.ErrorContains(t, err, tt.err.Error())
			} else {
				require.Equal(t, tt.err, err)
			}
		})
	}
}

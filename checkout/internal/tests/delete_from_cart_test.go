package tests

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	domain2 "route256/checkout/internal/domain"
	"route256/checkout/internal/models"
	"route256/checkout/internal/repository"
	repoMock "route256/checkout/internal/repository/mocks"
)

func TestDeleteFromCart(t *testing.T) {
	t.Parallel()
	type checkoutRepositoryMockFunc func(mc *minimock.Controller) repository.Repository

	type args struct {
		ctx    context.Context
		userID models.UserID
		sku    models.Sku
		count  models.Count
	}

	var (
		mc      = minimock.NewController(t)
		ctx     = context.Background()
		repoErr = errors.New("repo error")
		user    = models.UserID(gofakeit.Int64())
		sku     = models.Sku(gofakeit.Uint32())
		count   = models.Count(gofakeit.Uint32())
	)
	t.Cleanup(mc.Finish)

	tests := []struct {
		name                   string
		args                   args
		want                   error
		err                    error
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
			checkoutRepositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := repoMock.NewRepositoryMock(mc)
				mock.DeleteFromCartMock.Expect(ctx, user, sku, count).Return(nil)
				return mock
			},
		},
		{
			name: "delete from cart repo error",
			args: args{
				ctx:    ctx,
				userID: user,
				sku:    sku,
				count:  count,
			},
			want: repoErr,
			err:  repoErr,
			checkoutRepositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := repoMock.NewRepositoryMock(mc)
				mock.DeleteFromCartMock.Expect(ctx, user, sku, count).Return(repoErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			domain := domain2.NewMockBuisnessLogic(tt.checkoutRepositoryMock(mc))
			err := domain.DeleteFromCart(tt.args.ctx, tt.args.userID, tt.args.sku, tt.args.count)
			if tt.err != nil {
				require.ErrorContains(t, err, tt.err.Error())
			} else {
				require.Equal(t, tt.err, err)
			}
		})
	}
}

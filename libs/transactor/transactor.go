package transactor

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/multierr"
)

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate minimock -i TransactionManager -o ./mocks/ -s "_minimock.go"
//go:generate minimock -i DbClient -o ./mocks/ -s "_minimock.go"
//go:generate minimock -i github.com/jackc/pgx/v4.Tx -o ./mocks/tx_minimock.go

type QueryEngine interface {
	Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
}

type TransactionManager interface {
	RunRepeteableReade(ctx context.Context, fx func(ctxTX context.Context) error) error
}

type QueryEngineProvider interface {
	GetQueryEngine(ctx context.Context) QueryEngine // tx/pool
}

type DbClient interface {
	GetPool() *pgxpool.Pool
	BeginTx(context.Context, pgx.TxOptions) (pgx.Tx, error)
}

type manager struct {
	dbClient DbClient
}

func NewTransactor(dbClient DbClient) *manager {
	return &manager{
		dbClient: dbClient,
	}
}

type txkey string

const Key = txkey("tx")

func (tm *manager) RunRepeteableReade(ctx context.Context, fx func(ctxTX context.Context) error) error {
	tx, err := tm.dbClient.BeginTx(ctx,
		pgx.TxOptions{
			IsoLevel: pgx.RepeatableRead,
		})
	if err != nil {
		return err
	}

	if err := fx(context.WithValue(ctx, Key, tx)); err != nil {
		return multierr.Combine(err, tx.Rollback(ctx))
	}

	if err := tx.Commit(ctx); err != nil {
		return multierr.Combine(err, tx.Rollback(ctx))
	}

	return nil
}

func (tm *manager) GetQueryEngine(ctx context.Context) QueryEngine {
	tx, ok := ctx.Value(Key).(QueryEngine)
	if ok && tx != nil {
		return tx
	}

	return tm.dbClient.GetPool()
}

package db

import (
	"context"

	"github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgx"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type client struct {
	pool *pgxpool.Pool
	crdbpgx.Conn
}

func NewDbClient(pool *pgxpool.Pool) *client {
	return &client{
		pool: pool,
	}
}

func (c *client) GetPool() *pgxpool.Pool {
	return c.pool
}

func (c *client) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	return c.pool.BeginTx(ctx, txOptions)
}

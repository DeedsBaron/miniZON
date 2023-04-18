package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"route256/checkout/internal/config"
	"route256/checkout/pkg/db"
	"route256/libs/logger"
	"route256/libs/transactor"
)

const (
	MinConns        = 2
	MaxConns        = 10
	MaxConnIdleTime = time.Hour
	MaxConnLifetime = time.Hour

	loopCount = 3
	loopDelay = time.Second
)

// New create new db connection
func New(ctx context.Context, dbCfg *config.DbConfig) (transactor.DbClient, error) {
	var (
		err  error
		pool *pgxpool.Pool
	)

	pool, err = pgxpool.Connect(ctx, db.BuildDSN(dbCfg))
	if err != nil {
		return nil, errors.Wrap(err, "unable to create connection pool")
	}

	config := pool.Config()
	config.MaxConnIdleTime = MaxConnIdleTime
	config.MaxConnLifetime = MaxConnLifetime
	config.MinConns = MinConns
	config.MaxConns = MaxConns

	if err = db.ConnRetry(ctx, pool, loopCount, loopDelay); err != nil {
		return nil, errors.WithMessage(err, "cannot connect to database")
	}
	logger.Infof("successfully connected to database")

	return NewDbClient(pool), err
}

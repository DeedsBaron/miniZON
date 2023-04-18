package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"route256/checkout/internal/config"
)

func BuildDSN(cfg *config.DbConfig) string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.SslMode,
	)
}

func ConnRetry(ctx context.Context, pool *pgxpool.Pool, loopCount int, loopDelay time.Duration) error {
	var err error
	for i := 0; i < loopCount; i++ {
		err = pool.Ping(ctx)
		if err != nil {
			log.Printf("Waiting %v before looping again. Try N: %d, err: %v", loopDelay, i, err)
			time.Sleep(loopDelay)
		}
	}
	return err
}

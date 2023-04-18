package db

import (
	"context"
	"fmt"
	"log"
	"route256/loms/internal/config"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
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

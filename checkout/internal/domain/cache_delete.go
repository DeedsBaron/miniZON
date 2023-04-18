package domain

import (
	"context"

	"route256/checkout/internal/models"
)

func (b *Domain) CacheDelete(ctx context.Context, keys models.Keys) error {
	for _, key := range keys {
		err := b.cache.Delete(ctx, key)
		if err != nil {
			return err
		}
	}
	return nil
}

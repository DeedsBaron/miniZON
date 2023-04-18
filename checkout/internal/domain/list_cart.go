package domain

import (
	"context"
	"fmt"
	"sync"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"golang.org/x/time/rate"
	"route256/checkout/internal/config"
	"route256/checkout/internal/models"
	"route256/libs/logger"
	"route256/libs/workerpool"
)

func (b *Domain) ListCart(ctx context.Context, userID models.UserID) (*models.Cart, error) {
	childSpan, childCtx := opentracing.StartSpanFromContext(ctx, "domain: create_order")
	defer childSpan.Finish()

	cart, err := b.repo.ListCart(childCtx, userID)
	if err != nil {
		logger.Errorw("listing cart",
			"err", err.Error(),
			"component", "domain")
		return nil, errors.WithMessage(errors.WithMessage(err, "listCart"), "domain")
	}

	limiter := rate.NewLimiter(10, 10)

	taskSize := len(cart.Items)

	tasks := make([]workerpool.Task[models.SkuInfo, *models.ProductInfo], 0, taskSize)

	callback := func(sku models.SkuInfo) (*models.ProductInfo, error) {
		cachedProductInfo, ok := b.cache.Get(ctx, fmt.Sprintf("%v", sku.Sku))
		if ok {
			val, ok := cachedProductInfo.(models.ProductInfo)
			if ok {
				return &val, nil
			}
		}

		err = limiter.Wait(ctx)

		if err != nil {
			logger.Debugf("wait error", err)
			return nil, err
		}

		productInfo, err := b.productService.GetProduct(ctx, uint32(sku.Sku))
		if err != nil {
			logger.Errorw("product service",
				"err", err.Error(),
				"component", "domain")
			return nil, err
		}
		// индекс необходим для дальнейшего соотношения полученных данных из канала и корзиной
		productInfo.CartIndex = sku.CartIndex

		b.cache.Set(ctx, fmt.Sprintf("%v", sku.Sku), *productInfo)
		return productInfo, nil
	}

	// наполняем пул задач
	for i := 0; i < taskSize; i++ {
		tasks = append(tasks, workerpool.Task[models.SkuInfo, *models.ProductInfo]{
			Callback: callback,
			InArgs: models.SkuInfo{
				Sku:       cart.Items[i].Sku,
				CartIndex: i,
			},
		})
	}

	batchingPool, results := workerpool.NewPool[models.SkuInfo, *models.ProductInfo](ctx, config.Data.Workers.Amount)

	batchingPool.SubmitTasks(ctx, tasks)

	var wg sync.WaitGroup

	var totalPrice uint32

	// вычитываем результат работы воркеров в отдельной горутине
	wg.Add(1)
	go func() {
		defer wg.Done()

		for res := range results {
			cart.Items[res.CartIndex].Name = res.Name
			cart.Items[res.CartIndex].Price = res.Price
			totalPrice += res.Price
		}
	}()

	// в родительской рутине проверяем наличие ошибок, которые произошли в колбэке воркера
	for err := range batchingPool.GetErrorChan() {
		if err != nil {
			return nil, err
		}
	}

	wg.Wait()
	cart.TotalPrice = totalPrice
	return cart, nil
}

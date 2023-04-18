package cache

import (
	"container/list"
	"context"
	"fmt"
	"hash/fnv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"route256/libs/logger"
)

var (
	CacheHits = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "route256",
		Subsystem: "cache",
		Name:      "hits_total",
	})
	CacheErrorHits = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "route256",
		Subsystem: "cache",
		Name:      "error_hits_total",
	})
	CacheHistogramResponseTime = promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "route256",
		Subsystem: "cache",
		Name:      "histogram_cache_response_time_seconds",
		Buckets:   prometheus.ExponentialBuckets(0.0001, 2, 16),
	})
)

type Cache[K comparable] interface {
	Set(ctx context.Context, key K, value any) bool
	Get(ctx context.Context, key K) (any, bool)
	Delete(ctx context.Context, key K) error
}

type cacheLRU[K comparable] struct {
	duration time.Duration       // время жизни кэша
	Buckets  []*bucket[K]        // бакет - сущность, в которую будет происходить запись по ключу
	lruList  *list.List          // двусвязный список для работы с lru
	lruMap   map[K]*list.Element // мапа ключ - элемент листа - для удобства работы с листом
	capacity int                 // емкость (сколько можно записать прежде чем будет работать lru policy)
	size     int                 // реальное количество элементов в листе
	sync.RWMutex
}

type bucket[K comparable] struct {
	sync.RWMutex
	data map[K]value
}

type value struct {
	val       any       // значение, которое будет записано по ключу
	createdAt time.Time // время создание записи
}

// NewCacheWithTTL - создает инстанс кэша с поддержкой ttl. Запускает вотчера, который
// за временем жизни кэша в отдельной го рутине
func NewCacheWithTTL[K comparable](ctx context.Context, duration time.Duration, numBuckets, lruCapacity int) Cache[K] {
	buckets := make([]*bucket[K], numBuckets)

	for i := 0; i < numBuckets; i++ {
		buckets[i] = &bucket[K]{
			data: make(map[K]value),
		}
	}
	lruList := list.New()

	c := &cacheLRU[K]{
		duration: duration,
		Buckets:  buckets,
		lruList:  lruList,
		lruMap:   make(map[K]*list.Element),
		capacity: lruCapacity,
		size:     0,
	}

	// запускаем вотчера, который следит за ttl кэша в своей рутине
	go c.cleanUp(ctx)

	return c
}

// NewCacheWithoutTTL - создает инстанс кэша без поддержки ttl.
func NewCacheWithoutTTL[K comparable](ctx context.Context, numBuckets, lruCapacity int) Cache[K] {
	buckets := make([]*bucket[K], numBuckets)

	for i := 0; i < numBuckets; i++ {
		buckets[i] = &bucket[K]{
			data: make(map[K]value),
		}
	}
	lruList := list.New()

	c := &cacheLRU[K]{
		duration: time.Duration(0),
		Buckets:  buckets,
		lruList:  lruList,
		lruMap:   make(map[K]*list.Element),
		capacity: lruCapacity,
		size:     0,
	}

	return c
}

// NewCache - создает инстанс кэша с поддержкой ttl или без, в зависимости от duration
func NewCache[K comparable](ctx context.Context, duration time.Duration, numBuckets, lruCapacity int) Cache[K] {
	switch duration {
	case 0:
		return NewCacheWithoutTTL[K](ctx, numBuckets, lruCapacity)
	default:
		return NewCacheWithTTL[K](ctx, duration, numBuckets, lruCapacity)
	}
}

// getBucket - метод, позволяющий найти по ключу в каком бакете записано значение
func (c *cacheLRU[K]) getBucket(key K) (*bucket[K], error) {
	hash, err := c.fnv32(key)
	if err != nil {
		logger.Errorf("cant get bucket with key=%s", key)
		return nil, err
	}
	bucketIndex := hash % uint32(len(c.Buckets))
	return c.Buckets[bucketIndex], nil
}

func (c *cacheLRU[K]) Set(ctx context.Context, key K, val any) bool {
	b, err := c.getBucket(key)
	if err != nil {
		return false
	}

	b.Lock()
	defer b.Unlock()

	// проверяем нет ли такого значения по ключу и апдейтим, если есть
	_, found := b.data[key]
	if found {
		b.data[key] = value{
			val:       val,
			createdAt: time.Now(),
		}

		// Перемещаем вначало lru листа
		c.RLock()
		c.lruList.MoveToFront(c.lruMap[key])
		c.RUnlock()

		return true
	}
	// добавляем новое к/з в бакет
	b.data[key] = value{
		val:       val,
		createdAt: time.Now(),
	}
	c.size++

	c.Lock()
	// Добавляем вначало листа новую пару ключ/знач
	element := c.lruList.PushFront(key)
	c.lruMap[key] = element
	c.Unlock()

	// Проверяем что не вышли за границы емкости
	if c.size > c.capacity {
		c.Lock()
		element = c.lruList.Back()
		if element != nil {
			oldKey := c.lruList.Remove(element).(K)
			delete(c.lruMap, oldKey)

			b, err = c.getBucket(oldKey)
			if err != nil {
				return false
			}
			delete(b.data, key)

			logger.Debugw("removed due to lru",
				"key", oldKey)
			c.size--
			c.Unlock()
		}
	}

	logger.Debugw("set info to in-mem cacheLRU",
		"key", key,
		"value", val)
	return true
}

func (c *cacheLRU[K]) Get(ctx context.Context, key K) (any, bool) {
	b, err := c.getBucket(key)
	if err != nil {
		return nil, false
	}
	timeStart := time.Now()

	b.RLock()
	defer b.RUnlock()

	res, found := b.data[key]
	if !found {
		return nil, false
	}

	c.RLock()
	defer c.RUnlock()
	// перемещаем вначало нашего листа
	if val, ok := c.lruMap[key]; ok {
		c.lruList.MoveToFront(val)
	}

	logger.Debugw("get info from in-mem cacheLRU",
		"key", key,
		"value", res.val)

	CacheHits.Inc()
	elapsed := time.Since(timeStart)
	CacheHistogramResponseTime.Observe(elapsed.Seconds())

	return res.val, true
}

func (c *cacheLRU[K]) Delete(ctx context.Context, key K) error {
	b, err := c.getBucket(key)
	if err != nil {
		return err
	}
	b.Lock()
	defer b.Unlock()

	delete(b.data, key)
	return nil
}

func (c *cacheLRU[T]) cleanUp(ctx context.Context) {
	ticker := time.NewTicker(c.duration)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			logger.Debugf("context is canceled, cacheLRU ttl watcher is stopped")
			return
		case <-ticker.C:
			for _, b := range c.Buckets {
				b.Lock()
				for key, val := range b.data {
					if time.Since(val.createdAt) > c.duration {
						delete(b.data, key)
						logger.Debugf("deleted value due to ttl for key=%v", key)
					}
				}
				b.Unlock()
			}
		}
	}
}

// fnv32 - функция, генерирующая хэш по указаному входному ключу
func (c *cacheLRU[K]) fnv32(key K) (uint32, error) {
	h := fnv.New32a()
	_, err := h.Write([]byte(fmt.Sprintf("%v", key)))
	if err != nil {
		logger.Errorw("cant make hash",
			"key", key,
			"error", err)
		CacheErrorHits.Inc()
		return 0, err
	}
	return h.Sum32(), nil
}

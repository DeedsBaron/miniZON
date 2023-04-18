package rate

import (
	"sync"
	"time"
)

type Limiter struct {
	capacity     int           // максимальное кол/во токенов в бакете
	tokens       int           // текущее кол/воо токенов в бакете
	fillInterval time.Duration // интервал добавления токенов в бакет
	fillAmount   int           // кол/во токенов, которые добавляются в бакет через указанный интервал
	lastFillTime time.Time     // время, когда бакет последний раз пополнялся токенами
	lock         sync.Mutex    // мьютекс для избежания конкурентного изменения кол/ва токенов
}

// NewLimiter - создает лимитер, c указанной емкостью, временем рефила и количеством рефила
func NewLimiter(capacity int, fillInterval time.Duration, fillAmount int) *Limiter {
	return &Limiter{
		capacity:     capacity,
		tokens:       capacity,
		fillInterval: fillInterval,
		fillAmount:   fillAmount,
		lastFillTime: time.Now(),
	}
}

func (r *Limiter) Allow() bool {
	r.lock.Lock()
	defer r.lock.Unlock()

	// считаем сколько токенов должно быть сейчас в бакете
	now := time.Now()
	timeSinceLastFill := now.Sub(r.lastFillTime)
	expectedTokens := int(timeSinceLastFill / r.fillInterval * time.Duration(r.fillAmount))

	if expectedTokens > r.capacity {
		expectedTokens = r.capacity
	}

	// добавляем токены в бакет
	r.tokens += expectedTokens
	if r.tokens > r.capacity {
		r.tokens = r.capacity
	}
	r.lastFillTime = now

	// проверяем достаточно ли токенов для совершения действия
	if r.tokens > 0 {
		r.tokens--
		return true
	}
	return false
}

// Wait - метод ждет, до тех пор, пока не будет доступного токена, для совершения действия
func (r *Limiter) Wait() {
	for {
		if !r.Allow() {
			time.Sleep(r.fillInterval)
		} else {
			return
		}
	}
}

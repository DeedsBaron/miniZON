package workerpool

import (
	"context"
	"sync"
)

type Task[In, Out any] struct {
	Callback func(In) (Out, error)
	InArgs   In
}

type Pool[In, Out any] interface {
	SubmitTasks(context.Context, []Task[In, Out])
	Close()
	GetErrorChan() <-chan error
}

var _ Pool[any, any] = &p[any, any]{}

type p[In, Out any] struct {
	amountWorkers int

	wg sync.WaitGroup

	taskSource chan Task[In, Out]
	outSink    chan Out
	errorChan  chan error
}

func NewPool[In, Out any](ctx context.Context, amountWorkers int) (Pool[In, Out], <-chan Out) {
	pool := &p[In, Out]{
		amountWorkers: amountWorkers,
	}

	pool.bootstrap(ctx)

	return pool, pool.outSink
}

// GetErrorChan возвращает канал, в который записываются ошибки колбэка
func (pool *p[In, Out]) GetErrorChan() <-chan error {
	return pool.errorChan
}

// Close закрываем канал на выход, чтобы потребители могли выйти из := range
func (pool *p[In, Out]) Close() {
	// ждем пока все закончат работу
	pool.wg.Wait()

	// закрываем канал с ошибками, чтобы не было утечек
	close(pool.errorChan)

	close(pool.outSink)
}

// SubmitTasks не блокирующий метод, который отправляет таски в канал тасок и закрывает канал
func (pool *p[In, Out]) SubmitTasks(ctx context.Context, tasks []Task[In, Out]) {
	go func() {
		for _, task := range tasks {
			select {
			case <-ctx.Done():
				return

			case pool.taskSource <- task:
			}
		}
		// Больше задач не будет, закрываем канал, чтобы каждый воркер мог выйти из := range
		close(pool.taskSource)
	}()
}

func (pool *p[In, Out]) bootstrap(ctx context.Context) {
	pool.taskSource = make(chan Task[In, Out], pool.amountWorkers)
	pool.outSink = make(chan Out, pool.amountWorkers)
	pool.errorChan = make(chan error, pool.amountWorkers)

	ctxC, _ := context.WithCancel(ctx)

	for i := 0; i < pool.amountWorkers; i++ {
		pool.wg.Add(1)
		go func() {
			defer pool.wg.Done()
			worker(ctxC, pool.taskSource, pool.outSink, pool.errorChan)
		}()
	}
	//ждем окончания работы всех воркеров и закрываем каналы
	go func() {
		pool.Close()
	}()
}

func worker[In, Out any](
	ctx context.Context,
	taskSource <-chan Task[In, Out],
	resultSink chan<- Out,
	errorChan chan<- error,
) error {
	for task := range taskSource {
		select {
		case <-ctx.Done():
			return nil
		default:
			out, err := task.Callback(task.InArgs)
			if err != nil {
				errorChan <- err
				return err
			} else {
				resultSink <- out
			}
		}
	}
	return nil
}
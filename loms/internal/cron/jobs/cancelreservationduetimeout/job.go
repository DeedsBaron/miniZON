package cancelreservationduetimeout

import (
	"context"
	"time"

	"route256/libs/logger"
	"route256/libs/transactor"
	"route256/loms/internal/config"
	"route256/loms/internal/models"
	"route256/loms/internal/repository"

	"github.com/go-co-op/gocron"
)

type CancelReservationDueTimeout struct {
	cron *gocron.Scheduler
	repo repository.Repository
	tm   transactor.TransactionManager
}

func NewJob(repo *repository.LomsRepository,
	tm transactor.TransactionManager) *CancelReservationDueTimeout {
	s := gocron.NewScheduler(time.Local)

	return &CancelReservationDueTimeout{
		cron: s,
		repo: repo,
		tm:   tm,
	}
}

func (sh *CancelReservationDueTimeout) StartAsync(ctx context.Context) error {
	job, err := sh.cron.Cron(config.Data.CronJobs.CancelReservationDueTimeoutJob.Cron).Do(func() {
		err := sh.tm.RunRepeteableReade(ctx, func(ctxTX context.Context) error {
			orders, err := sh.repo.GetUnpayedOrdersWithinTimeout(ctx)
			if err != nil {
				return err
			}
			if len(orders) == 0 {
				return nil
			}

			err = sh.repo.UpdateOrderStatus(ctx, orders, models.StatusCanceled)
			if err != nil {
				return err
			}
			logger.Infow("reservations were canceled due timeout",
				"component", "CancelReservationDueTimeout job",
				"orderIDs", orders,
				"timeout", config.Data.CronJobs.CancelReservationDueTimeoutJob.OrderToBePayedTimeout)
			return nil
		})

		if err != nil {
			logger.Errorw("trying to make transaction",
				"component", "CancelReservationDueTimeout job",
				"err", err.Error())
		}
	})

	if err != nil {
		return err
	}
	if job.Error() != nil {
		return err
	}

	sh.cron.StartAsync()
	logger.Infof("cron job for checking orders status successfully started")
	return err
}

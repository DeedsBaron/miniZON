package readoutboxsend

import (
	"context"
	"time"

	"route256/libs/logger"
	"route256/loms/internal/config"
	"route256/loms/internal/models"
	"route256/loms/internal/producer"
	orderStatusChanges "route256/loms/internal/producer/order_status_changes"
	"route256/loms/internal/repository"

	"github.com/go-co-op/gocron"
)

type ReadOutboxSend struct {
	cron     *gocron.Scheduler
	repo     repository.Repository
	producer producer.KafkaProducer
}

func NewJob(repo *repository.LomsRepository,
	producer producer.KafkaProducer) *ReadOutboxSend {
	s := gocron.NewScheduler(time.Local)

	return &ReadOutboxSend{
		cron:     s,
		repo:     repo,
		producer: producer,
	}
}

func (j *ReadOutboxSend) StartAsync(ctx context.Context) error {
	job, err := j.cron.CronWithSeconds(config.Data.CronJobs.ReadOutBoxSendJob.Cron).Do(func() {
		outbox, err := j.repo.GetOutbox(ctx)
		if err != nil {
			logger.Errorw(err.Error(), "component", "ReadOutboxSend job")
		}

		successfullySent := 0
		for i, row := range outbox {
			key, value, err := orderStatusChanges.NewMessage(row)
			if err != nil {
				logger.Errorw("creating kafka message",
					"err", err.Error(),
					"component", "ReadOutboxSend job",
					"id", row.ID)
			}
			err = j.producer.SendMessage(ctx, key, value)
			if err != nil {
				logger.Errorw("sending message to kafka",
					"err", err.Error(),
					"component", "ReadOutboxSend job",
					"key", key,
					"value", value)
			}
			outbox[i].IsSent = true
			successfullySent++
		}
		batchToUpdate := make([]models.OutboxID, 0, successfullySent)
		for _, row := range outbox {
			if row.IsSent {
				batchToUpdate = append(batchToUpdate, row.ID)
			}
		}
		err = j.repo.UpdateOutbox(ctx, batchToUpdate)
		if err != nil {
			logger.Errorw("updating outbox",
				"err", err.Error(),
				"component", "ReadOutboxSend job")
		}
	})

	if err != nil {
		return err
	}
	if job.Error() != nil {
		return err
	}
	j.cron.StartAsync()
	logger.Infof("cron job for reading outbox and sending it to kafka successfully started")
	return err
}

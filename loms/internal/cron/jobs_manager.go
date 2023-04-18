package cron

import (
	"context"
)

type Job interface {
	StartAsync(ctx context.Context) error
}

type JobsManager struct {
	cancelReservationDueTimeoutJob Job
	readOutboxSendJob              Job
}

func NewJobsManager(cancelReservationDueTimeoutJob Job,
	readOutboxSendJob Job) *JobsManager {
	return &JobsManager{
		cancelReservationDueTimeoutJob: cancelReservationDueTimeoutJob,
		readOutboxSendJob:              readOutboxSendJob,
	}
}

func (m *JobsManager) StartAllJobs(ctx context.Context) error {
	err := m.cancelReservationDueTimeoutJob.StartAsync(ctx)
	if err != nil {
		return err
	}
	err = m.readOutboxSendJob.StartAsync(ctx)
	if err != nil {
		return err
	}
	return err
}

package di

import (
	"Jayleonc/gateway/internal/service"
	"Jayleonc/gateway/pkg/retry"
)

func InitRetryScheduler(demo *service.Demo) *retry.Scheduler {
	scheduler := retry.NewScheduler()
	scheduler.RegisterAsync(demo)
	return scheduler
}

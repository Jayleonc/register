package di

import (
	"Jayleonc/gateway/internal/client/log"
	"Jayleonc/gateway/internal/service"
	"Jayleonc/gateway/pkg/cronx"
	"github.com/robfig/cron/v3"
)

func InitJobs(l log.Logger, demo *service.Demo) *cronx.Cron {
	//builder := job.NewCronJobBuilder(l)

	var err error

	expr := cron.New(cron.WithSeconds())
	//_, err = expr.AddJob(job.EverySecond, builder.Build(demo))
	if err != nil {
		panic(err)
	}

	return &cronx.Cron{
		Cron: expr,
	}
}

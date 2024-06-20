package job

import (
	"Jayleonc/gateway/internal/client/log"
	"fmt"
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
)

// 常见的时间间隔表达式
const (
	// EverySecond 每秒钟运行一次
	EverySecond = "@every 1s"
	// EveryMinute 每分钟运行一次
	EveryMinute = "@every 1m"
	// EveryHour 每小时运行一次
	EveryHour = "@every 1h"
	// EveryDayMidnight 每天午夜运行一次
	EveryDayMidnight = "0 0 0 * * *"
	// EverySunday 每周日午夜运行一次
	EverySunday = "0 0 0 * * 0"
)

// 常见的 cron 表达式
const (
	// EveryMinuteAtZeroSecond 每分钟的第 0 秒运行一次
	EveryMinuteAtZeroSecond = "0 * * * * *"
	// EveryHourAtZeroMinute 每小时的第 0 分钟第 0 秒运行一次
	EveryHourAtZeroMinute = "0 0 * * * *"
	// EveryDayAtMidnight 每天的第 0 小时第 0 分钟第 0 秒运行一次
	EveryDayAtMidnight = "0 0 0 * * *"
	// EveryMonthFirstDay 每个月的第 1 天的第 0 小时第 0 分钟第 0 秒运行一次
	EveryMonthFirstDay = "0 0 0 1 * *"
	// EveryWeekSunday 每周日的第 0 小时第 0 分钟第 0 秒运行一次
	EveryWeekSunday = "0 0 0 * * 0"
	// FirstMondayOfMonth 每月的第一个星期一的0点0分0秒运行
	FirstMondayOfMonth = "0 0 0 * * 1#1"
)

var ModuleName = "CronJob"

type CronJobBuilder struct {
	l log.Logger
}

func NewCronJobBuilder(l log.Logger) *CronJobBuilder {
	return &CronJobBuilder{
		l: l,
	}
}

func (b *CronJobBuilder) Build(job Job) cron.Job {
	return cronJobAdapterFunc(func() {
		name := job.Name()
		requestId := uuid.New().String()
		b.l.Info(ModuleName, "Run", requestId, fmt.Sprintf("JobName: %s", name), nil)
		err := job.Run()
		if err != nil {
			b.l.Error(ModuleName, "Run", requestId, fmt.Sprintf("JobName: %s | Error message: %s", name, err.Error()), nil)
		}
	})
}

type cronJobAdapterFunc func()

func (c cronJobAdapterFunc) Run() {
	c()
}

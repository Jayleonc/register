package wire

import (
	"Jayleonc/gateway/internal/client/log"
	"Jayleonc/gateway/internal/events"
	"Jayleonc/gateway/pkg/cronx"
	"Jayleonc/gateway/pkg/ginx"
	"Jayleonc/gateway/pkg/retry"
)

type App struct {
	Web       *ginx.Server
	LogSender log.SenderI
	Scheduler *retry.Scheduler
	Cron      *cronx.Cron
	Consumers []events.Consumer
}

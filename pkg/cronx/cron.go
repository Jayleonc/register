package cronx

import "github.com/robfig/cron/v3"

type Cron struct {
	*cron.Cron
}

func (c *Cron) Start() {
	c.Cron.Start()
}

func (c *Cron) Stop() {
	c.Cron.Stop()
}

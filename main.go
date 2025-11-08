package main

import (
	"github.com/robfig/cron/v3"
	"github.com/yohannfra/cron-kuma-pusher/config"
	"github.com/yohannfra/cron-kuma-pusher/job"
)

func main() {
	c := cron.New()

	config.Init()
	config := config.GetConfig()

	for _, j := range config.Jobs {
		job.CreateJob(c, &j)
	}

	c.Start()

	select {}
}

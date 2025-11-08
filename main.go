package main

import (
	"flag"
	"log"

	"github.com/robfig/cron/v3"
	"github.com/yohannfra/cron-kuma-pusher/config"
	"github.com/yohannfra/cron-kuma-pusher/job"
)

func main() {
	configPath := flag.String("config", "", "Path to the configuration file (required)")
	flag.Parse()

	if *configPath == "" {
		log.Fatal("Error: --config flag is required. Please specify the path to your configuration file.")
	}

	c := cron.New()

	config.Init(*configPath)
	cfg := config.GetConfig()

	for _, j := range cfg.Jobs {
		job.CreateJob(c, &j)
	}

	c.Start()

	select {}
}

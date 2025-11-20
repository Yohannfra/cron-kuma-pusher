package main

import (
	"flag"
	"log"
	"os"

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

	c := cron.New(cron.WithSeconds())

	config.Init(*configPath)
	cfg := config.GetConfig()

	// create logs dir
	log.Printf("Creating logs directory '%s'", cfg.Logs.Dir)
	err := os.MkdirAll(cfg.Logs.Dir, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create logs dir: '%s'", cfg.Logs.Dir)
	}

	for _, j := range cfg.Jobs {
		job.CreateJob(c, &j)
	}

	c.Start()

	select {}
}

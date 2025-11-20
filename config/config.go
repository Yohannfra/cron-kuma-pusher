package config

import (
	"log"

	"github.com/spf13/viper"
)

type CronFormat string

const (
	FormatStandard CronFormat = "standard"
	FormatQuartz   CronFormat = "quartz"
)

type CronConfig struct {
	Format CronFormat
}

type UptimeKumaConfig struct {
	Enabled bool
	BaseUrl string
}

type LogsConfig struct {
	Enabled bool
	Dir     string
}

type Job struct {
	Name       string
	Expression string
	Command    string
	PushToken  string
}

type Configuration struct {
	Cron       CronConfig
	UptimeKuma UptimeKumaConfig
	Logs       LogsConfig
	Jobs       []Job
}

var c *Configuration

func validateConfig(c *Configuration) {
	// Check the cron config
	if c.Cron.Format != FormatStandard && c.Cron.Format != FormatQuartz {
		// if it's not defined set it to standard by default
		if c.Cron.Format == "" {
			c.Cron.Format = FormatStandard
		} else {
			log.Fatalf("Invalid cron format '%s'", c.Cron.Format)
		}
	}

	// Check if there is at least one job
	if len(c.Jobs) == 0 {
		log.Fatalf("Error: No jobs found in configuration")
	}

	// check if there is a duplicated name in jobs
	names := make(map[string]int)

	for _, job := range c.Jobs {
		names[job.Name]++
		if names[job.Name] > 1 {
			log.Fatalf("Found multiple jobs with name '%s'", job.Name)
		}
	}

	// Check if the KumaBaseUrl is present and valid
	if c.UptimeKuma.Enabled {
		if c.UptimeKuma.BaseUrl == "" {
			log.Fatal("Error: KumaBaseUrl is not set.")
		}
	}

	// Check if the LogsDir is present and set it's default value if not
	if c.Logs.Enabled {
		if c.Logs.Dir == "" {
			c.Logs.Dir = "./cron-kuma-pusher-logs"
		}
	}
}

func Init(configPath string) {
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigFile(configPath)

	v.SetDefault("jobs", []Job{})

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatalf("Error: Config file not found at %s. %v", configPath, err)
		} else {
			log.Fatalf("Fatal error reading config file: %v", err)
		}
	}

	log.Printf("Config file loaded: %v", v.ConfigFileUsed())

	if err := v.Unmarshal(&c); err != nil {
		log.Fatalf("Unable to unmarshal config into struct: %v", err)
	}

	validateConfig(c)

	// Cron config
	log.Println("Cron config:")
	log.Printf("- Format: %s", c.Cron.Format)

	// Uptime Kuma
	log.Println("Uptime Kuma config:")
	log.Printf("- Enabled: %v", c.UptimeKuma.Enabled)
	if c.UptimeKuma.Enabled {
		log.Printf("- Base url: %s", c.UptimeKuma.BaseUrl)
	}

	// Logs
	log.Println("Logs config:")
	log.Printf("- Enabled: %v", c.Logs.Enabled)
	if c.Logs.Enabled {
		log.Printf("- Directory %s", c.Logs.Dir)
	}

	// Jobs
	log.Println("Jobs:")
	log.Printf("- Count %d", len(c.Jobs))
	for _, job := range c.Jobs {
		log.Printf("- Name: %s:", job.Name)
		log.Printf("    Expression: %s", job.Expression)
		log.Printf("    Command: %s", job.Command)
		if c.UptimeKuma.Enabled {
			log.Printf("    Push token: %s", job.PushToken)
		}
	}
}

func GetConfig() *Configuration {
	return c
}

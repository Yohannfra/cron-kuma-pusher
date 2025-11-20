package config

import (
	"log"

	"github.com/spf13/viper"
)

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
	Jobs       []Job
	UptimeKuma UptimeKumaConfig
	Logs       LogsConfig
}

var c *Configuration

func validateConfig(c *Configuration) {
	// Check if there is at least one pinger
	if len(c.Jobs) == 0 {
		log.Fatalf("Error: No jobs found in configuration")
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

	// Uptime Kuma
	log.Printf("Uptime Kuma enabled is %v", c.UptimeKuma.Enabled)
	if c.UptimeKuma.Enabled {
		log.Printf("Uptime Kuma base url is %s", c.UptimeKuma.BaseUrl)
	}

	// Logs
	log.Printf("Logs enabled is %v", c.Logs.Enabled)
	if c.Logs.Enabled {
		log.Printf("Logs dir is %s", c.Logs.Dir)
	}

	// Jobs
	log.Printf("Loaded %d jobs from configuration", len(c.Jobs))
	for _, job := range c.Jobs {
		log.Println("--------------------------------")
		log.Printf("Name: %s", job.Name)
		log.Printf("Expression: %s", job.Expression)
		log.Printf("Command: %s", job.Command)
		log.Printf("Push token: %s", job.PushToken)
		log.Print("--------------------------------\n\n")
	}
}

func GetConfig() *Configuration {
	return c
}

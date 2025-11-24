package config

import (
	"log"
	"os"

	"github.com/goccy/go-yaml"
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
	BaseUrl string `yaml:"baseUrl"`
}

type LogsConfig struct {
	Enabled bool
	Dir     string
}

type Job struct {
	Name       string
	Expression string
	Workdir    string
	Command    string
	Env        []map[string]string
	EnvFile    string `yaml:"envFile"`
	PushToken  string `yaml:"pushToken"`
}

type Configuration struct {
	Cron       CronConfig
	UptimeKuma UptimeKumaConfig `yaml:"uptimeKuma"`
	Logs       LogsConfig
	Jobs       []Job
}

var c *Configuration = &Configuration{}

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
	fc, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	if err := yaml.Unmarshal(fc, c); err != nil {
		log.Fatalf("Unable to unmarshal config into struct: %v", err)
	}

	log.Printf("Config file loaded: %v", configPath)

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
		log.Printf("    Workdir: %s", job.Workdir)
		if len(job.Env) > 0 {
			log.Println("    Env:")
			for _, e := range job.Env {
				for k, v := range e {
					log.Printf("      %v=%v", k, v)
				}
			}
		}
		if job.EnvFile != "" {
			log.Printf("    envFile: %s", job.EnvFile)
		}
		log.Printf("    Command: %s", job.Command)
		if c.UptimeKuma.Enabled {
			log.Printf("    Push token: %s", job.PushToken)
		}
	}
}

func GetConfig() *Configuration {
	return c
}

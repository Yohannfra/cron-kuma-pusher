package config

import (
	"github.com/spf13/viper"
	"log"
)

type Job struct {
	Name       string
	Expression string
	Command    string
	PushToken  string
}

type Configuration struct {
	Jobs        []Job
	KumaBaseUrl string
}

var c *Configuration

func validateConfig(c *Configuration) {
	// Check if there is at least one pinger
	if len(c.Jobs) == 0 {
		log.Fatalf("Error: No jobs found in configuration")
	}

	// Check if the DiscordWebhookUrl is present and valid
	if c.KumaBaseUrl == "" {
		log.Fatal("Error: KumaBaseUrl is not set.")
	}
}

func Init() {
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigName("config")
	v.AddConfigPath(".")

	v.SetDefault("jobs", []Job{})

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatalf("Error: Config file not found. Make sure config.json is in the current directory. %v", err)
		} else {
			log.Fatalf("Fatal error reading config file: %v", err)
		}
	}

	log.Printf("Config file loaded: %v", v.ConfigFileUsed())

	if err := v.Unmarshal(&c); err != nil {
		log.Fatalf("Unable to unmarshal config into struct: %v", err)
	}

	validateConfig(c)

	log.Printf("Uptime Kuma base url is %s", c.KumaBaseUrl)
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

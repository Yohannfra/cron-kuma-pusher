package job

import (
	"log"
	"net/http"
	"net/url"

	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"github.com/yohannfra/cron-kuma-pusher/config"
	"github.com/yohannfra/cron-kuma-pusher/exec"
	"github.com/yohannfra/cron-kuma-pusher/utils"
)

func pushResultToKuma(pushToken, message string, exitCode int) {
	config := config.GetConfig()

	// "<KumaBaseUrl>/<token>?status=up&msg=OK"
	pushUrl := config.UptimeKuma.BaseUrl + "/" + pushToken + "?" +
		"status=" + utils.Ternary(exitCode == 0, "up", "down") +
		"&msg=" + url.QueryEscape(message)

	_, err := http.Get(pushUrl)
	if err != nil {
		log.Printf("Failed to push to kuma: %v", err)
	}
}

func CreateJob(c *cron.Cron, job *config.Job) {
	config := config.GetConfig()

	if job.EnvFile != "" {
		var loadedEnv map[string]string
		loadedEnv, err := godotenv.Read(job.EnvFile)

		if err != nil {
			log.Fatalf("Failed to load env file for job '%s'", job.Name)
		}

		// log.Printf("Loaded env file %s: %v", job.EnvFile, loadedEnv)

		job.Env = append(job.Env, loadedEnv)
		// log.Printf("Merged env for job %s: %v", job.Name, job.Env)
	}

	log.Printf("Creating job '%s'", job.Name)

	_, err := c.AddFunc(job.Expression, func() {
		stdout, stderr, exitCode, err := exec.Exec(job.Workdir, job.Env, job.Command)

		if err != nil {
			log.Printf("Error: failed to run command: %v\n", err)
			if config.Logs.Enabled {
				utils.AppendLog(job.Name, stdout, stderr, exitCode)
			}
			if config.UptimeKuma.Enabled {
				pushResultToKuma(job.PushToken, "Error: failed to run command", -1)
			}
			return
		}

		if config.Logs.Enabled {
			utils.AppendLog(job.Name, stdout, stderr, exitCode)
		}

		if exitCode == 0 {
			log.Printf("Job '%s' ran successfully", job.Name)
			if config.UptimeKuma.Enabled {
				pushResultToKuma(job.PushToken, "OK", 0)
			}
		} else {
			log.Printf("\n==== Job '%s' failed ====", job.Name)
			log.Printf("Command: %s\n", job.Command)
			log.Printf("Exit Code: %d\n", exitCode)
			log.Printf("Stdout:\n%s\n", stdout)
			log.Printf("Stderr:\n%s\n", stderr)
			log.Print("========================================\n\n")
			if config.UptimeKuma.Enabled {
				pushResultToKuma(job.PushToken, "KO", exitCode)
			}
		}
	})

	if err != nil {
		log.Fatalf("Failed to create job '%s' %v", job.Name, err.Error())
	}
}

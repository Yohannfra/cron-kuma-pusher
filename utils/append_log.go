package utils

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/yohannfra/cron-kuma-pusher/config"
)

func AppendLog(fp, stdout, stderr string, exitCode int) {
	cfg := config.GetConfig()

	f, err := os.OpenFile(path.Join(cfg.Logs.Dir, fp+".log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("error opening file: %v", err)
		return
	}
	defer f.Close()

	timestamp := time.Now().Format("2006-01-02 15:04:05")

	entry := fmt.Sprintf(
		"\n----- CRON RUN (%s) -----\nEXIT CODE: %d\nSTDOUT:\n%s\nSTDERR:\n%s\n---------------------------\n",
		timestamp,
		exitCode,
		stdout,
		stderr,
	)

	_, err = f.WriteString(entry + "\n")

	if err != nil {
		log.Printf("Failed to write to file: %v", err.Error())
	}
}

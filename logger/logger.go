package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"

	"go_demo_api/config"
)

var Log zerolog.Logger

func Init(cfg *config.Config) error {
	if err := os.MkdirAll("logs", 0755); err != nil {
		return err
	}

	logFile := "logs/" + cfg.ServiceName + ".log"
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	Log = zerolog.New(file).
		With().
		Timestamp().
		Str("service", cfg.ServiceName).
		Logger()

	go dailyRotate(cfg.ServiceName)

	return nil
}

func dailyRotate(serviceName string) {
	for {
		now := time.Now()
		nextRotate := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
		duration := nextRotate.Sub(now)

		time.Sleep(duration)

		oldFile := "logs/" + serviceName + ".log"
		newFile := "logs/" + now.Format("20060102") + serviceName + ".log"

		os.Rename(oldFile, newFile)

		file, err := os.OpenFile(oldFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			continue
		}

		Log = zerolog.New(file).
			With().
			Timestamp().
			Str("service", serviceName).
			Logger()
	}
}

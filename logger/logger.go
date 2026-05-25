package logger

import (
	"os"
	"path/filepath"
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
		Str("version", cfg.Version).
		Str("mode", cfg.EnvMode).
		Logger()

	go dailyRotate(cfg.ServiceName, cfg)

	return nil
}

func dailyRotate(serviceName string, cfg *config.Config) {
	for {
		now := time.Now()
		nextRotate := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
		duration := nextRotate.Sub(now)

		time.Sleep(duration)

		oldFile := filepath.Join("logs/", serviceName+".log")
		newFile := filepath.Join("logs/", now.Format("20060102")+"_"+serviceName+".log")

		os.Rename(oldFile, newFile)

		file, err := os.OpenFile(oldFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			continue
		}

		Log = zerolog.New(file).
			With().
			Timestamp().
			Str("service", serviceName).
			Str("version", cfg.Version).
			Str("mode", cfg.EnvMode).
			Logger()
	}
}

package jobs

import (
	"log"
)

type CleanupJob struct{}

func (j CleanupJob) Name() string {
	return "cleanup-job"
}

func (j CleanupJob) Schedule() string {
	return "0 */1 * * *" // every hour
}

func (j CleanupJob) Run() {
	log.Println("Cleaning temporary files...")
}

package jobs

import (
	"log"
)

type EmailJob struct{}

func (j EmailJob) Name() string {
	return "email-job"
}

func (j EmailJob) Schedule() string {
	return "*/5 * * * *" // every 5 min
}

func (j EmailJob) Run() {
	log.Println("Sending emails...")
}

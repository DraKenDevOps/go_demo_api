package scheduler

import (
	"log"

	"github.com/robfig/cron/v3"
)

type Job interface {
	Name() string
	Schedule() string
	Run()
}

type Scheduler struct {
	cron *cron.Cron
}

func New() *Scheduler {
	return &Scheduler{
		cron: cron.New(),
	}
}

func (s *Scheduler) Register(job Job) error {
	_, err := s.cron.AddFunc(job.Schedule(), func() {
		log.Printf("Running job: %s", job.Name())
		job.Run()
	})

	return err
}

func (s *Scheduler) Start() {
	s.cron.Start()
}

func (s *Scheduler) Stop() {
	ctx := s.cron.Stop()
	<-ctx.Done()
}

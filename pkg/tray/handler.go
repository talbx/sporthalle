package tray

import (
	"log/slog"
	"os"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/talbx/sporthalle/pkg/collect"
	"github.com/talbx/sporthalle/pkg/eval"
	"github.com/talbx/sporthalle/pkg/types"
)

func OnReady() {
	s, err := gocron.NewScheduler()
	if err != nil {
		slog.Default().Error("could not create scheduler")
		os.Exit(1)
	}

	service := NewSystrayService()
	_, err = s.NewJob(
		gocron.DurationJob(
			30*time.Second,
		),
		gocron.NewTask(
			executionFunc, s, service,
		),
		gocron.WithStartAt(gocron.WithStartImmediately()))
	if err != nil {
		slog.Default().Error("could not run job")
		slog.Default().Error(err.Error())
		os.Exit(1)
	}
	s.Start()
	types.LOGGER.Info("Successfully started Sporthalle tray widget!")

}

func executionFunc(s gocron.Scheduler, svc *SystrayService) {
	events := collect.Run()
	todayEvent, nextEvent := eval.UpcomingEvents(events)
	// there should be only exactly one job
	nextRun, _ := s.Jobs()[0].NextRun()
	context := ExecutionContext{
		CurrentEvent: todayEvent,
		NextEvent:    nextEvent,
		NextRun:      nextRun,
	}

	svc.Tell(context)
	go svc.Quit()
}

type ExecutionContext struct {
	CurrentEvent *types.Event
	NextEvent    types.Event
	NextRun      time.Time
}

func OnExit() {
	// clean up here
}

package tray

import (
	"github.com/go-co-op/gocron/v2"
	"github.com/talbx/sporthalle/pkg/eval"
	"github.com/talbx/sporthalle/pkg/run"
	"github.com/talbx/sporthalle/pkg/types"
	"log/slog"
	"os"
	"time"
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
	events := run.Run()
	todayEvent := eval.IsEvent(events)
	// there should be only exactly one job
	nextRun, _ := s.Jobs()[0].NextRun()
	context := ExecutionContext{
		Event:   todayEvent,
		NextRun: nextRun,
	}

	svc.Tell(context)
	go svc.Quit()
}

type ExecutionContext struct {
	Event   *types.Event
	NextRun time.Time
}

func OnExit() {
	// clean up here
}

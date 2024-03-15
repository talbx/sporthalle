package tray

import (
	"fmt"
	"github.com/getlantern/systray"
	"github.com/go-co-op/gocron/v2"
	"github.com/talbx/sporthalle/pkg/eval"
	"github.com/talbx/sporthalle/pkg/run"
	"github.com/talbx/sporthalle/pkg/types"
	"log/slog"
	"os"
	"time"
)

type Tray struct {
	EventIndicator   *systray.MenuItem
	RefreshIndicator *systray.MenuItem
	Quitter          *systray.MenuItem
}

func OnReady() {
	s, err := gocron.NewScheduler()
	if err != nil {
		slog.Default().Error("could not create scheduler")
		os.Exit(1)
	}

	tray := NewTray()

	_, err = s.NewJob(
		gocron.DurationJob(
			30*time.Second,
		),
		gocron.NewTask(
			executionFunc, s, tray,
		),
		gocron.WithStartAt(gocron.WithStartImmediately()))
	if err != nil {
		slog.Default().Error("could not run job")
		os.Exit(1)
	}
	s.Start()
	types.LOGGER.Info("Successfully started Sporthalle tray widget!")

}

func NewTray() Tray {
	return Tray{
		EventIndicator:   systray.AddMenuItem("", ""),
		RefreshIndicator: systray.AddMenuItem("", ""),
		Quitter:          systray.AddMenuItem("Quit", "Quit the whole app"),
	}
}

func executionFunc(s gocron.Scheduler, tray Tray) {
	events := run.Run()
	todayEvent := eval.IsEvent(events)
	// there should be only exactly one job
	nextRun, _ := s.Jobs()[0].NextRun()
	context := ExecutionContext{
		Event:   todayEvent,
		NextRun: nextRun,
		Tray:    tray,
	}

	determine(context)
	go quit(context)
}

type ExecutionContext struct {
	Event   *types.Event
	NextRun time.Time
	Tray    Tray
}

func determine(ec ExecutionContext) {

	go provideRefreshMenuItem(ec)

	if nil == ec.Event {
		handleNoEvent(ec)
		return
	}
	handleEventToday(ec)
}

func provideRefreshMenuItem(ec ExecutionContext) {
	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {
		timeRemaining := getTimeRemaining(ec.NextRun)

		if timeRemaining <= 0 {
			break
		}
		ec.Tray.RefreshIndicator.SetTitle(fmt.Sprintf("Next update in: %d\n", timeRemaining))
	}
}

func handleEventToday(ec ExecutionContext) {
	systray.SetTitle("⚠️")
	s := fmt.Sprintf("%s @ %s %s", ec.Event.Name, ec.Event.Date.Format("02.01.2006"), ec.Event.Start)
	ec.Tray.EventIndicator.SetTitle(s)
	systray.SetIcon(Red)
}

func handleNoEvent(ec ExecutionContext) {
	ec.Tray.EventIndicator.SetTitle("No Event today")
	systray.SetIcon(Green)
}

func OnExit() {
	// clean up here
}

func quit(ec ExecutionContext) {
	<-ec.Tray.Quitter.ClickedCh
	systray.Quit()
}

func getTimeRemaining(t time.Time) int {
	currentTime := time.Now()
	difference := t.Sub(currentTime)

	total := int(difference.Seconds())
	seconds := total % 60

	return seconds
}

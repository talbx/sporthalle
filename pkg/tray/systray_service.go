package tray

import (
	"fmt"
	"github.com/getlantern/systray"
	"time"
)

type SystrayService struct {
	eventIndicator   *systray.MenuItem
	refreshIndicator *systray.MenuItem
	quitButton       *systray.MenuItem
	NextRun          time.Time
}

func NewSystrayService() *SystrayService {
	return &SystrayService{
		eventIndicator:   systray.AddMenuItem("", "Today's event"),
		refreshIndicator: systray.AddMenuItem("", "When will the event refresh"),
		quitButton:       systray.AddMenuItem("Quit", "Quits Sporthalle"),
	}
}

func (s *SystrayService) Tell(executionCtx ExecutionContext) {
	go s.provideRefreshMenuItem(executionCtx)

	if nil == executionCtx.Event {
		s.handleNoEvent(executionCtx)
		return
	}
	s.handleEventToday(executionCtx)
}

func (s *SystrayService) provideRefreshMenuItem(ec ExecutionContext) {
	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {
		timeRemaining := getTimeRemaining(ec.NextRun)

		if timeRemaining <= 0 {
			break
		}
		s.refreshIndicator.SetTitle(fmt.Sprintf("Next update in: %d\n", timeRemaining))
	}
}

func (s *SystrayService) handleEventToday(ec ExecutionContext) {
	systray.SetTitle("⚠️")
	text := fmt.Sprintf("%s @ %s %s", ec.Event.Name, ec.Event.Date.Format("02.01.2006"), ec.Event.Start)
	s.eventIndicator.SetTitle(text)
	systray.SetIcon(Red)
}

func (s *SystrayService) handleNoEvent(ec ExecutionContext) {
	s.eventIndicator.SetTitle("No Event today")
	systray.SetIcon(Green)
}

func getTimeRemaining(t time.Time) int {
	currentTime := time.Now()
	difference := t.Sub(currentTime)

	total := int(difference.Seconds())
	seconds := total % 60

	return seconds
}

func (s *SystrayService) Quit() {
	<-s.quitButton.ClickedCh
	systray.Quit()
}

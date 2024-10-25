package tray

import (
	"fmt"
	"github.com/getlantern/systray"
	"github.com/talbx/sporthalle/pkg/core/types"
	"log"
	"os/exec"
	"runtime"
	"time"
)

type SystrayService struct {
	eventIndicator     *systray.MenuItem
	nextEventIndicator *systray.MenuItem
	refreshIndicator   *systray.MenuItem
	quitButton         *systray.MenuItem
	NextRun            time.Time
}

func NewSystrayService() *SystrayService {
	return &SystrayService{
		eventIndicator:     systray.AddMenuItem("", "Today's event"),
		nextEventIndicator: systray.AddMenuItem("", "Next Event"),
		refreshIndicator:   systray.AddMenuItem("", "When will the event refresh"),
		quitButton:         systray.AddMenuItem("Quit", "Quits Sporthalle"),
	}
}

func (s *SystrayService) Tell(executionCtx ExecutionContext) {
	go s.provideRefreshMenuItem(executionCtx)
	go s.handleEventClick(s.eventIndicator.ClickedCh)
	go s.handleEventClick(s.nextEventIndicator.ClickedCh)

	if nil == executionCtx.CurrentEvent {
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
		s.refreshIndicator.SetTitle(fmt.Sprintf("Refresh in: %d\n", timeRemaining))
	}
}

func (s *SystrayService) handleEventToday(ec ExecutionContext) {
	systray.SetTitle("⚠️")
	currentText := s.createEventMessage(*ec.CurrentEvent)
	nextText := s.createEventMessage(ec.NextEvent)
	s.eventIndicator.SetTitle(currentText)
	s.nextEventIndicator.SetTitle(nextText)
	systray.SetIcon(Red)
}

func (s *SystrayService) handleNoEvent(ec ExecutionContext) {
	s.eventIndicator.SetTitle("No Event today!")
	s.nextEventIndicator.SetTitle("Next Event: " + s.createEventMessage(ec.NextEvent))
	systray.SetIcon(Green)
}

func (s *SystrayService) createEventMessage(event types.Event) string {
	return fmt.Sprintf("%s @ %s %s", event.Name, event.Date.Format("02.01.2006"), event.Start)
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

func (s *SystrayService) handleEventClick(ch chan struct{}) {
	<-ch
	err := open("https://ssl.webpack.de/termine.sporthallehamburg.de/pr/clipper.php")
	if err != nil {
		log.Default().Print(err)
	}
}

func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

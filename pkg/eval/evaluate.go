package eval

import (
	"github.com/talbx/sporthalle/pkg/types"
	"time"
)

func IsEvent(events []types.Event) *types.Event {
	now := time.Now()
	for _, event := range events {
		if event.Date.Year() == now.Year() && event.Date.YearDay() == now.YearDay() {
			types.LOGGER.Info("There is an event today! will send pushover message", "event", event.Name, "date", event.Date)
			return &event
		}
	}
	return nil
}

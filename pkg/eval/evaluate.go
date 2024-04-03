package eval

import (
	"slices"
	"time"

	"github.com/talbx/sporthalle/pkg/types"
)

func UpcomingEvents(events []types.Event) (today *types.Event, next types.Event) {
	now := time.Now()

	slices.SortFunc(events, func(a, b types.Event) int {
		if a.Date.Before(b.Date) {
			return -1
		} else if a.Date.YearDay() == b.Date.YearDay() {
			return 0
		}
		return 1
	})

	var todaysEvent *types.Event = nil
	for _, event := range events {
		if todaysEvent != nil {
			return todaysEvent, event
		}
		if event.Date.Year() == now.Year() && event.Date.YearDay() == now.YearDay() {
			types.LOGGER.Info("There is an event today", "event", event.Name, "date", event.Date)
			todaysEvent = &event
		}
	}
	return nil, events[0]
}

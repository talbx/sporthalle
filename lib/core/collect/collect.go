package collect

import (
	"fmt"
	"github.com/talbx/sporthalle/lib/core/types"
	"log/slog"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

var Sporthalle = "https://ssl.webpack.de/termine.sporthallehamburg.de/pr/clipper.php"
var lastEntry = "Beginn: first"

type Collector struct {
	log *slog.Logger
}

func NewCollector(log *slog.Logger) Collector {
	return Collector{log: log}
}

func (c *Collector) Run() []types.Event {
	events, err := c.collect()
	if err != nil {
		c.log.Error(err.Error())
		os.Exit(1)
	}

	c.log.Info("sucessfully created event list", "event size", len(events))
	return events
}

func (m *Collector) collect() ([]types.Event, error) {
	events := make([]string, 0)
	dates := make([]string, 0)
	starts := make([]string, 0)
	c := colly.NewCollector()

	c.OnHTML(".rahmen_radius_l", func(e *colly.HTMLElement) {
		events = append(events, e.Text)
	})

	c.OnHTML("div", func(e *colly.HTMLElement) {
		if e.Attr("style") == "margin-left:4px;margin-bottom:2px;" {
			dateMatch, _ := regexp.MatchString("^(Mo|Di|Mi|Do|Fr|Sa|So) [0-9][0-9]*.[0-9][0-9]*.202[0-9]", e.Text)
			if dateMatch {
				if !strings.HasPrefix(lastEntry, "Einlass") && !strings.HasPrefix(lastEntry, "Beginn") {
					m.log.Debug(fmt.Sprintf("The last entry for %v does not have a starting time. will add n/a as starting time", lastEntry))
					starts = append(starts, "Einlass n/a")
					lastEntry = "Einlass: n/a"
				}
				dates = append(dates, e.Text)
				lastEntry = e.Text
			}
		}

		if e.Attr("style") == "" {
			startMatch, _ := regexp.MatchString("^(Einlass:|Beginn:) [0-9][0-9]:[0-9][0-9]", e.Text)
			if startMatch {
				if !strings.HasPrefix(lastEntry, "Einlass") && !strings.HasPrefix(lastEntry, "Beginn") {
					starts = append(starts, e.Text)
					lastEntry = e.Text
				}
			}
		}
	})

	err := c.Visit(Sporthalle)
	if err != nil {
		return nil, err
	}
	eventys := make([]types.Event, 0)

	// fix last element incase it has no start time
	if len(events) == len(dates) && len(dates) > len(starts) {
		starts = append(starts, "Einlass n/a")
	}
	m.log.Info("calculation concluded.", "event size", len(events), "date size", len(dates), "time size", len(starts))

	for i, event := range events {
		date := dates[i]
		start := starts[i]

		realDate := strings.Split(date, " ")
		parse, err := time.Parse("02.01.2006", realDate[1])
		if err != nil {
			return nil, err
		}
		ev := types.Event{
			Name: event,
			// January 2, 2006
			Date:  parse,
			Start: start,
		}
		eventys = append(eventys, ev)
	}
	return eventys, nil

}

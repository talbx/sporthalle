package types

import (
	"github.com/gocolly/colly"
	"regexp"
	"strings"
	"time"
)

type SporthallenCollector struct{}

var Sporthalle = "https://ssl.webpack.de/termine.sporthallehamburg.de/pr/clipper.php"
var eventList = make([]Event, 0)

func (sc SporthallenCollector) Run() ([]Event, error) {
	events := make([]string, 0)
	dates := make([]string, 0)
	starts := make([]string, 0)
	begins := make([]string, 0)
	c := colly.NewCollector()

	// Find and visit all links
	c.OnHTML(".rahmen_radius_l", func(e *colly.HTMLElement) {
		eventLen := len(events)
		dateLen := len(dates)
		startLen := len(starts)

		validate(eventLen, dateLen, startLen, dates, events, starts)

		events = append(events, e.Text)

	})

	c.OnHTML("div", func(e *colly.HTMLElement) {
		dateMatch, _ := regexp.MatchString("^(Mo|Di|Mi|Do|Fr|Sa|So) [0-9][0-9]*.[0-9][0-9]*.202[0-9]", e.Text)
		startMatch, _ := regexp.MatchString("^(Einlass:) [0-9][0-9]:[0-9][0-9]", e.Text)

		if dateMatch {
			dates = append(dates, e.Text)
			return
		}
		if startMatch {
			split := strings.Split(e.Text, "Einlass: ")
			starts = append(starts, split[1])
			return
		}

		if beginMatch, _ := regexp.MatchString("^(Beginn:) [0-9][0-9]:[0-9][0-9]", e.Text); beginMatch {
			split := strings.Split(e.Text, "Beginn: ")
			begins = append(begins, split[1])
			if len(dates) != len(starts) {
				starts = append(starts, split[1])
			}

			return
		}

		validate(len(events), len(dates), len(starts), dates, events, starts)

		//fmt.Println("no date and no start match found!", e.Text)

	})

	err := c.Visit(Sporthalle)
	if err != nil {
		return nil, err
	}

	LOGGER.Info("calculation concluded.", "event size", len(events), "date size", len(dates), "time size", len(starts))

	eventys := make([]Event, 0)
	for i, event := range events {
		date := dates[i]
		start := starts[i]

		realDate := strings.Split(date, " ")
		parse, err := time.Parse("02.01.2006", realDate[1])
		if err != nil {
			return nil, err
		}
		ev := Event{
			Name: event,
			// January 2, 2006
			Date:  parse,
			Start: start,
		}

		eventys = append(eventys, ev)

	}
	return eventys, nil

}

func validate(eventLen int, dateLen int, startLen int, dates []string, events []string, starts []string) {
	if eventLen == dateLen && dateLen == startLen && eventLen != 0 {
		d := dates[eventLen-1]
		split := strings.Split(d, " ")
		parse, err := time.Parse("02.01.2006", split[1])
		if err != nil {
			panic(err)
		}
		eventList = append(eventList, Event{
			Name:  events[eventLen-1],
			Date:  parse,
			Start: starts[eventLen-1],
		})
	}
}

package types

import (
	"log/slog"
	"os"
	"time"
)

type Collector interface {
	Run() ([]Event, error)
}

var LOGGER = slog.New(slog.NewJSONHandler(os.Stdout, nil))

type Notifier interface {
	Notify(events []Event) (string, error)
}

type Event struct {
	Name  string
	Date  time.Time
	Start string
}

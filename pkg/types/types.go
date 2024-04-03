package types

import (
	"log/slog"
	"os"
	"time"
)

var LOGGER = slog.New(slog.NewJSONHandler(os.Stdout, nil))
type Event struct {
	Name  string
	Date  time.Time
	Start string
}

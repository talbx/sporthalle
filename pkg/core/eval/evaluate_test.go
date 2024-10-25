package eval

import (
	"github.com/stretchr/testify/assert"
	"github.com/talbx/sporthalle/pkg/core/types"
	"testing"
	"time"
)

func TestIsEventTrue(t *testing.T) {
	events := []types.Event{{
		Name:  "Architects",
		Start: "18:00",
		Date:  time.Now(),
	}, {
		Name:  "Loathe",
		Start: "18:00",
		Date:  time.Now().AddDate(0, 1, 0),
	}}

	today, _ := UpcomingEvents(events)
	assert.NotNil(t, today)
	assert.Equal(t, *today, events[0])
}

func TestIsEventNo(t *testing.T) {
	events := []types.Event{{
		Name:  "Helene Fischer",
		Start: "18:00",
		Date:  time.Now().AddDate(1, 0, 0),
	}, {
		Name:  "Loathe",
		Start: "18:00",
		Date:  time.Now().AddDate(0, 1, 0),
	}}

	today, _ := UpcomingEvents(events)
	assert.Nil(t, today)
}

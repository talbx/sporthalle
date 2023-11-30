package eval

import (
	"github.com/stretchr/testify/assert"
	"github.com/talbx/sporthalle/pkg/types"
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

	event := IsEvent(events)
	assert.NotNil(t, event)
	assert.Equal(t, *event, events[0])
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

	event := IsEvent(events)
	assert.Nil(t, event)
}

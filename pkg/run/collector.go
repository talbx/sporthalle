package run

import (
	"github.com/talbx/sporthalle/pkg/types"
	"os"
)

var notifyIt = false

func Run() []types.Event {
	collector := types.CollectorFactory{}.Create("")
	run := runCollector(collector)
	return run
}

func runCollector(co types.Collector) []types.Event {
	events, err := co.Run()
	if err != nil {
		types.LOGGER.Error(err.Error())
		os.Exit(1)
	}
	types.LOGGER.Info("sucessfully created event list", "event size", len(events))

	if notifyIt {
		var n types.Notifier = PushoverNotifier{}

		notify, err := n.Notify(events)
		if err != nil {
			types.LOGGER.Error(err.Error())
			os.Exit(1)
		}
		types.LOGGER.Info("successfully notified.", "id", notify)
	}
	return events
}

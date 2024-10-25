package fn

import (
	"context"
	"errors"
	"fmt"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/talbx/sporthalle/pkg/core/collect"
	"github.com/talbx/sporthalle/pkg/notify"
	"log/slog"
	"os"
)

func init() {
	functions.CloudEvent("sporthalle", receive)
}

type MessagePublishedData struct {
	Message PubSubMessage
}

// PubSubMessage is the payload of a Pub/Sub event.
// See the documentation for more details:
// https://cloud.google.com/pubsub/docs/reference/rest/v1/PubsubMessage
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// receive consumes a CloudEvent message and extracts the Pub/Sub message.
func receive(ctx context.Context, e event.Event) error {
	var msg MessagePublishedData
	if err := e.DataAs(&msg); err != nil {
		return fmt.Errorf("event.DataAs: %v", err)
	}

	name := string(msg.Message.Data) // Automatically decoded from base64.
	if name == "sporthalle" {
		if ctx.Value("test").(string) == "yes" {
			return errors.New("test scenario ends here")
		}
		log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
		collector := collect.NewCollector(log)
		events := collector.Run()
		n := notify.NewPushoverNotifier(log)

		no, err := n.Notify(events)
		if err != nil {
			log.Error(err.Error())
			os.Exit(1)
		}
		log.Info("successfully notified.", "id", no)
	}
	return fmt.Errorf("could not do anything with msg %s", name)
}

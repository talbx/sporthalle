package run

import (
	"fmt"
	"github.com/gregdel/pushover"
	"github.com/talbx/sporthalle/pkg/eval"
	"github.com/talbx/sporthalle/pkg/types"
	"gopkg.in/yaml.v3"
	"os"
)

type PushoverNotifier struct {
}

func (p PushoverNotifier) Notify(events []types.Event) (string, error) {
	event := eval.IsEvent(events)
	if event != nil {
		types.LOGGER.Info("There is an event today! will send pushover message", "event", event.Name, "date", event.Date)
		return p.sendMessage(*event)
	}
	return "", nil
}

type PushoverConfig struct {
	ApiToken  string `yaml:"apiToken"`
	UserToken string `yaml:"userToken"`
}

func (p PushoverNotifier) sendMessage(event types.Event) (string, error) {
	file, err := os.ReadFile("config.yaml")
	if err != nil {
		return "", err
	}
	var conf PushoverConfig
	err = yaml.Unmarshal(file, &conf)
	if err != nil {
		return "", err
	}
	message := &pushover.Message{
		Message:  fmt.Sprintf("at %s %s %s starts", event.Date.Format("02.01.2006"), event.Start, event.Name),
		Title:    fmt.Sprintf("Heute: %s in der Sporthalle (%s)", event.Name, event.Start),
		Priority: pushover.PriorityNormal,
		Sound:    pushover.SoundSiren,
		URLTitle: "sporthalle termin übersicht",
		URL:      types.Sporthalle,
	}
	recipient := pushover.NewRecipient(conf.UserToken)
	response, err := pushover.New(conf.ApiToken).SendMessage(message, recipient)
	if err != nil {
		return "", err
	}
	return response.ID, nil
}

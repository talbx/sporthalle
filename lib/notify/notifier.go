package notify

import (
	"fmt"
	"github.com/talbx/sporthalle/lib/core/eval"
	"github.com/talbx/sporthalle/lib/core/types"
	"log/slog"
	"os"
	"time"

	"github.com/gregdel/pushover"
	"gopkg.in/yaml.v3"
)

type Notifier interface {
	Notify(events []types.Event) (string, error)
}
type PushoverNotifier struct {
	log *slog.Logger
}

func NewPushoverNotifier(log *slog.Logger) PushoverNotifier {
	return PushoverNotifier{log}
}

func (p PushoverNotifier) Notify(events []types.Event) (string, error) {
	event, next := eval.UpcomingEvents(events)
	if event != nil {
		slog.Default().Info("There is an event today! will send pushover message", "event", event.Name, "date", event.Date)
		return p.notifyToday(*event)
	}
	return p.notifyNext(next)
}

type PushoverConfig struct {
	ApiToken  string `yaml:"apiToken"`
	UserToken string `yaml:"userToken"`
}

func (p *PushoverNotifier) notifyNext(event types.Event) (string, error) {
	ev := event.Date.YearDay()
	now := time.Now().YearDay()
	daysLeft := ev - now
	p.log.Info(fmt.Sprintf("%v days left until next event", daysLeft))
	message := &pushover.Message{
		Message:  fmt.Sprintf("%s at %s %s", event.Name, event.Date.Format("02.01.2006"), event.Start),
		Title:    fmt.Sprintf("Next Event in %v days!", daysLeft),
		Priority: pushover.PriorityNormal,
		Sound:    pushover.SoundBike,
		URLTitle: "sporthalle termin übersicht",
		URL:      "https://ssl.webpack.de/termine.sporthallehamburg.de/pr/clipper.php",
	}
	return p.sendMessage(message)
}

func (p *PushoverNotifier) sendMessage(message *pushover.Message) (string, error) {
	file, err := os.ReadFile("config.yaml")
	if err != nil {
		return "", err
	}
	var conf PushoverConfig
	err = yaml.Unmarshal(file, &conf)
	if err != nil {
		return "", err
	}
	recipient := pushover.NewRecipient(conf.UserToken)
	response, err := pushover.New(conf.ApiToken).SendMessage(message, recipient)
	if err != nil {
		return "", err
	}
	return response.ID, nil
}

func (p *PushoverNotifier) notifyToday(event types.Event) (string, error) {
	message := &pushover.Message{
		Message:  fmt.Sprintf("at %s %s %s starts", event.Date.Format("02.01.2006"), event.Start, event.Name),
		Title:    fmt.Sprintf("Heute: %s in der Sporthalle (%s)", event.Name, event.Start),
		Priority: pushover.PriorityNormal,
		Sound:    pushover.SoundSiren,
		URLTitle: "sporthalle termin übersicht",
		URL:      "https://ssl.webpack.de/termine.sporthallehamburg.de/pr/clipper.php",
	}
	return p.sendMessage(message)
}

package tray

import (
	"fmt"
	"github.com/getlantern/systray"
	"github.com/talbx/sporthalle/pkg/eval"
	"github.com/talbx/sporthalle/pkg/run"
	"github.com/talbx/sporthalle/pkg/types"
)

func OnReady() {
	events := run.Run()
	todayEvent := eval.IsEvent(events)
	var mQuit *systray.MenuItem
	mQuit = determine(todayEvent, mQuit)

	// Sets the icon of a menu item. Only available on Mac and Windows.
	go quit(mQuit)
	types.LOGGER.Info("Successfully started Sporthalle tray widget!")
	types.LOGGER.Info("Today's event", "event", todayEvent)
}

func determine(todayEvent *types.Event, mQuit *systray.MenuItem) *systray.MenuItem {
	if nil == todayEvent {
		return handleNoEvent(mQuit)
	}
	return handleEventToday(todayEvent, mQuit)
}

func handleEventToday(todayEvent *types.Event, mQuit *systray.MenuItem) *systray.MenuItem {
	systray.SetTitle("⚠️")
	s := fmt.Sprintf("%s @ %s %s", todayEvent.Name, todayEvent.Date.Format("02.01.2006"), todayEvent.Start)
	systray.AddMenuItem(s, "Today's Event")
	mQuit = systray.AddMenuItem("Exit", "Quit the whole app")
	systray.SetIcon(Red)
	return mQuit
}

func handleNoEvent(mQuit *systray.MenuItem) *systray.MenuItem {
	systray.SetTitle("")
	systray.AddMenuItem("No Event today", "lucky you!")
	mQuit = systray.AddMenuItem("Exit", "Quit the whole app")
	systray.SetIcon(Green)
	return mQuit
}

func OnExit() {
	// clean up here
}

func quit(mQuit *systray.MenuItem) {
	types.LOGGER.Info("Exit requested...")
	<-mQuit.ClickedCh
	systray.Quit()
}

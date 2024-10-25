package main

import (
	"github.com/getlantern/systray"
	"github.com/talbx/sporthalle/lib/tray"
)

func main() {
	systray.Run(tray.OnReady, tray.OnExit)
}

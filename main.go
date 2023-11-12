package main

import (
	"github.com/getlantern/systray"
	"github.com/talbx/sporthalle/pkg/tray"
)

func main() {
	systray.Run(tray.OnReady, tray.OnExit)
}

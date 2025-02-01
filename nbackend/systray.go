package backend

import (
	"github.com/getlantern/systray"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"log"
	goruntime "runtime"
)

func (a *App) setupSystemTray() {
	systray.Register(
		func() {
			// OnReady
			systray.SetIcon(getIcon())
			systray.SetTitle("GPT Snapper")
			systray.SetTooltip("AI Screenshot Assistant")

			mOpen := systray.AddMenuItem("Open Settings", "Open configuration")
			mQuit := systray.AddMenuItem("Quit", "Exit the application")

			go func() {
				for {
					select {
					case <-mOpen.ClickedCh:
						runtime.WindowShow(a.ctx)
					case <-mQuit.ClickedCh:
						systray.Quit()
						runtime.Quit(a.ctx)
					}
				}
			}()
		},
		func() {
			// OnExit
			log.Println("System tray exited")
		},
	)
}

func getIcon() []byte {
	// Return embedded icon bytes for your OS
	if goruntime.GOOS == "windows" {
		return []byte{ /* Windows .ico bytes */ }
	}
	return []byte{ /* Mac/Linux .png bytes */ }
}

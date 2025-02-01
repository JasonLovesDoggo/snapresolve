package backend

import (
	"github.com/emersion/go-autostart"
)

type AutoStart struct {
	app *autostart.App
}

func NewAutoStart() *AutoStart {
	return &AutoStart{
		app: &autostart.App{
			Name:        "GPTSnapper",
			DisplayName: "GPT Screenshot Assistant",
			Exec:        []string{"/usr/local/bin/gpt-snapper"},
		},
	}
}

func (a *AutoStart) Enable() error   { return a.app.Enable() }
func (a *AutoStart) Disable() error  { return a.app.Disable() }
func (a *AutoStart) IsEnabled() bool { return a.app.IsEnabled() }

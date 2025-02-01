package services

//
//import (
//	"github.com/emersion/go-autostart"
//)
//
//type AutoStart struct {
//	app *autostart.App
//}
//
//func NewAutoStart() *AutoStart {
//	return &AutoStart{
//		app: &autostart.App{
//			Name:        "SnapResolve",
//			DisplayName: "SnapResolve Daemon",
//			Exec:        []string{"/usr/local/bin/snapresolve"},
//		},
//	}
//}
//
//func (a *AutoStart) Enable() error   { return a.app.Enable() }
//func (a *AutoStart) Disable() error  { return a.app.Disable() }
//func (a *AutoStart) IsEnabled() bool { return a.app.IsEnabled() }

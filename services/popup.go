package services

import (
	"github.com/ncruces/zenity"
)

type UIService struct{}

func NewUIService() *UIService {
	return &UIService{}
}

func (u *UIService) ShowResult(result string) error {
	// First show a progress dialog
	progress, err := zenity.Progress(
		zenity.Title("Analyzing Screenshot"),
		zenity.NoCancel(),
	)
	if err != nil {
		return err
	}

	progress.Close()

	// Then show the result
	err = zenity.Info(result,
		zenity.Title("Screenshot Analysis"),
		zenity.Width(400),
		zenity.InfoIcon,
	)
	return err
}

func (u *UIService) ShowError(err string) error {
	err2 := zenity.Error(err,
		zenity.Title("Error"),
		zenity.Width(400),
	)
	return err2
}

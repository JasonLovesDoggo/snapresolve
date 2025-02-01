package backend

import (
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) showResponsePopup(content string) {
	runtime.WindowExecJS(a.ctx, `
		document.getElementById('response-popup').classList.remove('hidden');
		document.getElementById('response-content').innerText = `+content+`;
	`)
}

func (a *App) showErrorPopup(message string) {
	runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Type:    runtime.ErrorDialog,
		Title:   "GPT Snapper Error",
		Message: message,
	})
}

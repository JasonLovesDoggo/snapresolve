package backend

import (
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"os"
	"time"
)

type App struct {
	ctx       context.Context
	config    *Config
	hotkey    *HotkeyManager
	autoStart *AutoStart
	llm       *LLMService
}

func NewApp() *App {
	return &App{
		config:    InitConfig(),
		autoStart: NewAutoStart(),
		hotkey:    NewHotkeyManager(),
		llm:       NewLLMService(),
	}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	a.setupSystemTray()

	// Apply config
	a.hotkey.Register(a.config.Hotkey)
	if a.config.AutoStart {
		err := a.autoStart.Enable()
		if err != nil {
			runtime.EventsEmit(ctx, "error", "Failed to enable auto-start")
		}
	}

	// Listen for hotkey events
	runtime.EventsOn(ctx, "hotkey-pressed", func(...interface{}) {
		a.handleScreenshot()
	})
}

func (a *App) handleScreenshot() {
	// Capture screenshot
	imagePath, err := CaptureScreenshot()
	if err != nil {
		a.showErrorPopup("Failed to capture screenshot: " + err.Error())
		return
	}
	defer os.Remove(imagePath) // Clean up

	// Get analysis from LLM
	response, err := a.llm.AnalyzeScreenshot(imagePath, a.config.Prompt)
	if err != nil {
		a.showErrorPopup("AI analysis failed: " + err.Error())
		return
	}
	for _, chunk := range splitResponse(response) {
		runtime.EventsEmit(a.ctx, "gpt-response", chunk)
		time.Sleep(50 * time.Millisecond)
	}

	// Show response
	//a.showResponsePopup(response)
}

func splitResponse(text string) []string {
	// todo: stream later
	return []string{text}
}

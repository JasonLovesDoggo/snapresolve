package main

import (
	"context"
	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"github.com/jasonlovesdoggo/snapresolve/services"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"os"
)

type App struct {
	ctx        context.Context
	config     *services.Config
	screenshot *services.ScreenshotService
	llm        *services.LLMService
}

func NewApp() *App {
	return &App{}
}

func (a *App) handleScreenshot() {
	// Capture screenshot
	imgPath, err := a.screenshot.CaptureCurrentScreen()
	if err != nil {
		runtime.LogError(a.ctx, "Failed to capture screenshot: "+err.Error())
		return
	}

	// Analyze with LLM
	result, err := a.llm.Analyze(imgPath)
	if err != nil {
		runtime.LogError(a.ctx, "Failed to analyze screenshot: "+err.Error())
		return
	}

	// Emit event to frontend with result
	runtime.EventsEmit(a.ctx, "analysis-result", result)

	// Clean up the temporary screenshot file
	os.Remove(imgPath)
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	runtime.WindowHide(ctx)

	cfg, err := services.LoadConfig()
	if err != nil {
		return
	}
	a.config = cfg

	a.screenshot = services.NewScreenshotService(cfg.TempDir)
	a.llm = services.NewLLMService(cfg.OpenAIKey, cfg.Prompt)
}

func (a *App) onReady() {
	systray.SetIcon(icon.Data) // You'll need to provide icon data
	systray.SetTitle("Snapresolve")
	systray.SetTooltip("AI Screenshot Analysis")

	mCapture := systray.AddMenuItem("Take Screenshot", "Analyze screen with AI")
	mSettings := systray.AddMenuItem("Settings", "Open settings")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Exit application")

	go func() {
		for {
			select {
			case <-mCapture.ClickedCh:
				a.handleScreenshot()
			case <-mSettings.ClickedCh:
				runtime.WindowShow(a.ctx)
			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}

func (a *App) onExit() {
	a.screenshot.CleanupTempFiles()
}

// AnalyzeScreenshot is exported for frontend use
func (a *App) AnalyzeScreenshot() (string, error) {
	img, err := a.screenshot.CaptureCurrentScreen()
	if err != nil {
		return "", err
	}

	return a.llm.Analyze(img)
}

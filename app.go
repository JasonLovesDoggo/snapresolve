package main

import (
	"context"
	"fmt"
	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"github.com/jasonlovesdoggo/snapresolve/services"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/viper"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"os"
)

type App struct {
	ctx        context.Context
	config     *services.Config
	screenshot *services.ScreenshotService
	llm        *services.LLMService
	ui         *services.UIService
	hotkey     *services.HotkeyService
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Load configuration
	cfg, err := services.LoadConfig()
	if err != nil {
		return
	}
	a.config = cfg

	// Initialize services
	a.screenshot = services.NewScreenshotService(cfg.TempDir)
	a.ui = services.NewUIService()

	// Initialize LLM service
	provider := services.Provider(cfg.Provider)
	apiKey := cfg.GeminiKey
	if provider == services.ProviderOpenAI {
		apiKey = cfg.OpenAIKey
	}

	llmService, err := services.NewLLMService(provider, apiKey, cfg.Prompt)
	if err != nil {
		a.ui.ShowError("Failed to initialize LLM service: " + err.Error())
		return
	}
	a.llm = llmService

	// Initialize hotkey service
	hotkeyService, err := services.NewHotkeyService(cfg.HotkeyCapture, a.handleScreenshot)
	if err != nil {
		a.ui.ShowError("Failed to register hotkey: " + err.Error())
		return
	}
	a.hotkey = hotkeyService
}

func (a *App) handleScreenshot() {

	fmt.Println("Taking screenshot...")

	// Capture screenshot
	imgPath, err := a.screenshot.CaptureCurrentScreen()
	if err != nil {
		a.ui.ShowError("Failed to capture screenshot: " + err.Error())
		return
	}
	defer os.Remove(imgPath) // Clean up the temporary file

	// Analyze with LLM
	result, err := a.llm.Analyze(imgPath)
	if err != nil {
		a.ui.ShowError("Failed to analyze screenshot: " + err.Error())
		return
	}

	// Show result
	if err := a.ui.ShowResult(result); err != nil {
		a.ui.ShowError("Failed to show result: " + err.Error())
	}
}

func (a *App) onExit() {
	if a.hotkey != nil {
		a.hotkey.Stop()
	}
	if err := a.screenshot.CleanupTempFiles(); err != nil {
		a.ui.ShowError("Failed to cleanup temp files: " + err.Error())
	}
}

func (a *App) Systray() {
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
				a.OpenSettings()
			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}

func (a *App) OpenSettings() {
	//configPath := getConfigPath()
	configFile := viper.ConfigFileUsed()
	runtime.LogInfof(a.ctx, "Opening config file: %s", configFile)

	err := open.Start(configFile)
	if err != nil {
		panic(err)
	}
}

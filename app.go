package main

import (
	"context"
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

	cfg, err := services.LoadConfig()
	if err != nil {
		runtime.LogError(ctx, "Failed to load config: "+err.Error())
		return
	}
	a.config = cfg

	a.screenshot = services.NewScreenshotService(cfg.TempDir)

	// Determine which API key to use based on provider
	provider := services.Provider(cfg.Provider)
	var apiKey string
	switch provider {
	case services.ProviderOpenAI:
		if cfg.OpenAIKey == "" {
			runtime.LogFatal(ctx, "OpenAI key is required")
			runtime.Quit(ctx)
		}
		apiKey = cfg.OpenAIKey
	case services.ProviderGemini:
		if cfg.GeminiKey == "" {
			runtime.LogFatal(ctx, "Gemini key is required")
			runtime.Quit(ctx)
		}
		apiKey = cfg.GeminiKey

	}

	// Initialize LLM service
	llmService, err := services.NewLLMService(provider, apiKey, cfg.Prompt)
	if err != nil {
		runtime.LogError(ctx, "Failed to initialize LLM service: "+err.Error())
		return
	}
	a.llm = llmService
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
				a.OpenSettings()
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

func (a *App) OpenSettings() {
	//configPath := getConfigPath()
	configFile := viper.ConfigFileUsed()
	runtime.LogInfof(a.ctx, "Opening config file: %s", configFile)

	err := open.Start(configFile)
	if err != nil {
		panic(err)
	}
}

package services

import (
	"errors"
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

type Config struct {
	OpenAIKey     string // OpenAI API key
	HotkeyCapture string // Hotkey to capture a screenshot
	TempDir       string // Directory to store the temporary screenshot files in (will be created if it doesn't exist)
	Prompt        string // Prompt for the user on what to do with the screenshot
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")

	configPath := getConfigPath()
	viper.AddConfigPath(configPath)

	viper.SetDefault("hotkey_capture", "Alt+Shift+S")
	viper.SetDefault("temp_dir", filepath.Join(configPath, "temp"))
	viper.SetDefault("prompt", "Please analyze this screenshot and say how to fix what is in it.")

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			err = viper.SafeWriteConfig()
			if err != nil {
				return nil, err
			}
		}
	}

	config := &Config{
		OpenAIKey:     viper.GetString("openai_key"),
		HotkeyCapture: viper.GetString("hotkey_capture"),
		TempDir:       viper.GetString("temp_dir"),
		Prompt:        viper.GetString("prompt"),
	}

	return config, nil
}

func getConfigPath() string {
	configDir, _ := os.UserConfigDir()

	configDir = filepath.Join(configDir, "snapresolve")

	os.MkdirAll(configDir, 0755)
	return configDir
}

func openSettings() {
	//configPath := getConfigPath()
	configFile := viper.ConfigFileUsed()

	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "notepad"
		args = []string{configFile}
	case "darwin":
		cmd = "open"
		args = []string{configFile}
	default: // linux
		cmd = "xdg-open"
		args = []string{configFile}
	}

	exec.Command(cmd, args...).Start()
}

package services

import (
	"errors"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

type Config struct {
	OpenAIKey     string `mapstructure:"openai_key"`
	GeminiKey     string `mapstructure:"gemini_key"`
	HotkeyCapture string `mapstructure:"hotkey_capture"`
	TempDir       string `mapstructure:"temp_dir"`
	Prompt        string `mapstructure:"prompt"`
	Provider      string `mapstructure:"provider"` // "openai" or "gemini"
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")

	configPath := getConfigPath()
	viper.AddConfigPath(configPath)

	// Set defaults
	viper.SetDefault("hotkey_capture", "Alt+Shift+S")
	viper.SetDefault("temp_dir", filepath.Join(configPath, "temp"))
	viper.SetDefault("prompt", "Please analyze this screenshot and say how to fix what is in it.")
	viper.SetDefault("provider", "gemini")

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			err = viper.SafeWriteConfig()
			if err != nil {
				return nil, err
			}
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func getConfigPath() string {
	configDir, _ := os.UserConfigDir()

	configDir = filepath.Join(configDir, "snapresolve")

	os.MkdirAll(configDir, 0755)
	return configDir
}

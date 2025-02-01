package backend

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	APIKey    string
	Hotkey    string
	LLM       string
	AutoStart bool
	Prompt    string
}

func InitConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(getConfigPath())

	if err := viper.ReadInConfig(); err != nil {
		viper.SetDefault("hotkey", "Ctrl+Shift+S")
		viper.SetDefault("llm", "gpt-4-vision")
		viper.SetDefault("auto_start", false)
		viper.SetDefault("prompt", "How do I fix the issue I am having in this screenshot?")
	}

	return &Config{
		APIKey:    viper.GetString("api_key"),
		Hotkey:    viper.GetString("hotkey"),
		LLM:       viper.GetString("llm"),
		AutoStart: viper.GetBool("auto_start"),
		Prompt:    viper.GetString("prompt"),
	}
}

func (c *Config) Save() error {
	viper.Set("api_key", c.APIKey)
	viper.Set("hotkey", c.Hotkey)
	viper.Set("llm", c.LLM)
	viper.Set("auto_start", c.AutoStart)
	viper.Set("prompt", c.Prompt)
	return viper.WriteConfig()
}

func getConfigPath() string {
	dir, _ := os.UserConfigDir()
	return filepath.Join(dir, "gpt-snapper")
}

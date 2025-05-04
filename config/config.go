package config

import (
	"encoding/json"
	"os"
)

// Config contains bot settings
type Config struct {
	Token           string `json:"token"`            // Discord bot token
	Prefix          string `json:"prefix"`           // Command prefix
	GeminiAPIKey    string `json:"gemini_api_key"`   // API key for Gemini
	GrokAPIKey      string `json:"grok_api_key"`     // API key for Grok
	ChatGPTAPIKey   string `json:"chatgpt_api_key"`  // API key for ChatGPT
	QwenAPIKey      string `json:"qwen_api_key"`     // API key for Qwen
	ClaudeAPIKey    string `json:"claude_api_key"`   // API key for Claude
	DefaultAI       string `json:"default_ai"`       // Default AI provider
	ReportThreshold int    `json:"report_threshold"` // Report threshold for auto-ban
	AdminRoleID     string `json:"admin_role_id"`    // Administrator role ID
	ModRoleID       string `json:"mod_role_id"`      // Moderator role ID
	DefaultLanguage string `json:"default_language"` // Default bot language (ru, en, uk, de, zh)
	BotName         string `json:"bot_name"`         // Discord bot name
}

// Load loads configuration from config.json file
func Load() (*Config, error) {
	file, err := os.Open("config.json")
	if err != nil {
		// If the file doesn't exist, create default configuration
		if os.IsNotExist(err) {
			defaultConfig := &Config{
				Prefix:          "/",
				ReportThreshold: 3,
				DefaultLanguage: "ru",
			}

			// Create file with default configuration
			if err := SaveConfig(defaultConfig); err != nil {
				return nil, err
			}

			return defaultConfig, nil
		}
		return nil, err
	}
	defer file.Close()

	config := &Config{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(config)
	return config, err
}

// SaveConfig saves configuration to file
func SaveConfig(config *Config) error {
	file, err := os.Create("config.json")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(config)
}

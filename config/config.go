package config

import (
	"encoding/json"
	"os"
)

// WebInterfaceConfig содержит настройки веб-интерфейса
type WebInterfaceConfig struct {
	Enabled  bool   `json:"enabled"`   // Включен ли веб-интерфейс
	Host     string `json:"host"`      // Хост для веб-интерфейса
	Port     int    `json:"port"`      // Основной порт для веб-интерфейса
	AltPorts []int  `json:"alt_ports"` // Альтернативные порты для веб-интерфейса
}

// Config contains bot settings
type Config struct {
	Token           string             `json:"token"`            // Discord bot token
	Prefix          string             `json:"prefix"`           // Command prefix
	GeminiAPIKey    string             `json:"gemini_api_key"`   // API key for Gemini
	GrokAPIKey      string             `json:"grok_api_key"`     // API key for Grok
	ChatGPTAPIKey   string             `json:"chatgpt_api_key"`  // API key for ChatGPT
	QwenAPIKey      string             `json:"qwen_api_key"`     // API key for Qwen
	ClaudeAPIKey    string             `json:"claude_api_key"`   // API key for Claude
	DefaultAI       string             `json:"default_ai"`       // Default AI provider
	ReportThreshold int                `json:"report_threshold"` // Report threshold for auto-ban
	AdminRoleID     string             `json:"admin_role_id"`    // Administrator role ID
	ModRoleID       string             `json:"mod_role_id"`      // Moderator role ID
	DefaultLanguage string             `json:"default_language"` // Default bot language (ru, en, uk, de, zh)
	BotName         string             `json:"bot_name"`         // Discord bot name
	WebInterface    WebInterfaceConfig `json:"web_interface"`    // Web interface settings
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
				WebInterface: WebInterfaceConfig{
					Enabled:  true,
					Host:     "localhost",
					Port:     8080,
					AltPorts: []int{3000, 8000},
				},
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
	// Установка значений по умолчанию для веб-интерфейса
	config.WebInterface.Enabled = true
	config.WebInterface.Host = "localhost"
	config.WebInterface.Port = 8080
	config.WebInterface.AltPorts = []int{3000, 8000}

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

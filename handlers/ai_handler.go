package handlers

import (
	"fmt"
	"strings"

	"discord-bot/ai"
	"discord-bot/localization"

	"github.com/bwmarrin/discordgo"
)

// Глобальные переменные для хранения команд приложения
var (
	aiCommands []*discordgo.ApplicationCommand
	aiHandlers map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)
)

// InitAICommands инициализирует команды AI для Discord
func InitAICommands(s *discordgo.Session) error {
	// Определяем обработчики команд
	aiHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"ai":      handleAIInteraction,
		"gemini":  handleAIModelInteraction,
		"grok":    handleAIModelInteraction,
		"chatgpt": handleAIModelInteraction,
		"qwen":    handleAIModelInteraction,
		"claude":  handleAIModelInteraction,
	}

	// Определяем команды приложения
	aiCommands = []*discordgo.ApplicationCommand{
		{
			Name:        "ai",
			Description: "Задать вопрос искусственному интеллекту (используя модель по умолчанию)",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "запрос",
					Description: "Ваш вопрос или запрос к AI",
					Required:    true,
				},
			},
		},
		{
			Name:        "gemini",
			Description: "Задать вопрос Gemini AI",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "запрос",
					Description: "Ваш вопрос или запрос к Gemini",
					Required:    true,
				},
			},
		},
		{
			Name:        "grok",
			Description: "Задать вопрос Grok AI",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "запрос",
					Description: "Ваш вопрос или запрос к Grok",
					Required:    true,
				},
			},
		},
		{
			Name:        "chatgpt",
			Description: "Задать вопрос ChatGPT",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "запрос",
					Description: "Ваш вопрос или запрос к ChatGPT",
					Required:    true,
				},
			},
		},
		{
			Name:        "qwen",
			Description: "Задать вопрос Qwen AI",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "запрос",
					Description: "Ваш вопрос или запрос к Qwen",
					Required:    true,
				},
			},
		},
		{
			Name:        "claude",
			Description: "Задать вопрос Claude AI",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "запрос",
					Description: "Ваш вопрос или запрос к Claude",
					Required:    true,
				},
			},
		},
	}

	// Регистрируем обработчик интеракций
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type == discordgo.InteractionApplicationCommand {
			if h, ok := aiHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
		}
	})

	// Регистрируем команды в Discord
	for _, cmd := range aiCommands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, "", cmd)
		if err != nil {
			return fmt.Errorf("не удалось создать команду %s: %w", cmd.Name, err)
		}
	}

	return nil
}

// HandleAICommand обрабатывает текстовые команды AI (для обратной совместимости)
func HandleAICommand(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) == 0 {
		s.ChannelMessageSend(m.ChannelID, localization.GetText("ai_usage", cfg.Prefix))
		return
	}

	// Проверяем, указана ли модель AI
	var modelName string
	var prompt string

	if len(args) > 1 && isAIModel(args[0]) {
		modelName = args[0]
		prompt = strings.Join(args[1:], " ")
	} else {
		modelName = ""
		prompt = strings.Join(args, " ")
	}

	// Отправляем сообщение о том, что запрос обрабатывается
	s.ChannelMessageSend(m.ChannelID, localization.GetText("ai_processing"))

	// Получаем ответ от AI
	response, err := generateAIResponse(modelName, prompt)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, localization.GetText("ai_error", err.Error()))
		return
	}

	// Отправляем ответ
	sendAIResponse(s, m.ChannelID, response)
}

// handleAIInteraction обрабатывает слеш-команду /ai
func handleAIInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	prompt := options[0].StringValue()

	// Отправляем сообщение о том, что запрос обрабатывается
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	// Получаем ответ от AI по умолчанию
	response, err := generateAIResponse("", prompt)
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, &discordgo.WebhookParams{
			Content: localization.GetText("ai_error", err.Error()),
		})
		return
	}

	// Отправляем ответ
	sendAIInteractionResponse(s, i, response)
}

// handleAIModelInteraction обрабатывает слеш-команды для конкретных моделей AI
func handleAIModelInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	modelName := i.ApplicationCommandData().Name
	options := i.ApplicationCommandData().Options
	prompt := options[0].StringValue()

	// Отправляем сообщение о том, что запрос обрабатывается
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	// Получаем ответ от указанной модели AI
	response, err := generateAIResponse(modelName, prompt)
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, &discordgo.WebhookParams{
			Content: localization.GetText("ai_error", err.Error()),
		})
		return
	}

	// Отправляем ответ
	sendAIInteractionResponse(s, i, response)
}

// generateAIResponse генерирует ответ от указанной модели AI
func generateAIResponse(modelName string, prompt string) (string, error) {
	var response string
	var err error

	if modelName == "" {
		// Используем провайдер по умолчанию
		response, err = ai.GenerateResponse(prompt)
	} else {
		// Используем указанный провайдер
		provider, err := ai.GetProvider(modelName)
		if err != nil {
			return "", err
		}
		response, err = provider.GenerateResponse(prompt)
	}

	return response, err
}

// sendAIResponse отправляет ответ от AI в текстовый канал
func sendAIResponse(s *discordgo.Session, channelID string, response string) {
	// Если ответ слишком длинный, разбиваем его на части
	if len(response) > 2000 {
		chunks := splitMessage(response, 2000)
		for _, chunk := range chunks {
			s.ChannelMessageSend(channelID, chunk)
		}
	} else {
		s.ChannelMessageSend(channelID, response)
	}
}

// sendAIInteractionResponse отправляет ответ от AI в ответ на интеракцию
func sendAIInteractionResponse(s *discordgo.Session, i *discordgo.InteractionCreate, response string) {
	// Если ответ слишком длинный, разбиваем его на части
	if len(response) > 2000 {
		chunks := splitMessage(response, 2000)
		for index, chunk := range chunks {
			if index == 0 {
				// Первый чанк отправляем как основной ответ
				s.FollowupMessageCreate(i.Interaction, &discordgo.WebhookParams{
					Content: chunk,
				})
			} else {
				// Остальные чанки отправляем как дополнительные сообщения
				s.FollowupMessageCreate(i.Interaction, &discordgo.WebhookParams{
					Content: chunk,
				})
			}
		}
	} else {
		s.FollowupMessageCreate(i.Interaction, &discordgo.WebhookParams{
			Content: response,
		})
	}
}

// isAIModel проверяет, является ли строка названием модели AI
func isAIModel(name string) bool {
	_, exists := ai.AvailableProviders[name]
	return exists
}

// splitMessage разбивает сообщение на части указанной длины
func splitMessage(message string, chunkSize int) []string {
	var chunks []string
	for len(message) > 0 {
		if len(message) <= chunkSize {
			chunks = append(chunks, message)
			break
		}

		// Ищем последний пробел в пределах chunkSize
		lastSpace := strings.LastIndex(message[:chunkSize], " ")
		if lastSpace == -1 {
			// Если пробел не найден, просто разбиваем по размеру чанка
			lastSpace = chunkSize
		}

		chunks = append(chunks, message[:lastSpace])
		message = message[lastSpace+1:]
	}
	return chunks
}

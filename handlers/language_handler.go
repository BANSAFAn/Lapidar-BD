package handlers

import (
	"discord-bot/localization"

	"github.com/bwmarrin/discordgo"
)

// HandleLanguageCommand обрабатывает команду смены языка
func HandleLanguageCommand(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	// Проверяем, что пользователь указал язык
	if len(args) < 1 {
		s.ChannelMessageSend(m.ChannelID, localization.GetText("language_usage", cfg.Prefix))
		return
	}

	// Получаем код языка
	langCode := args[0]

	// Устанавливаем новый язык
	if success := localization.SetLanguage(langCode); !success {
		s.ChannelMessageSend(m.ChannelID, localization.GetText("language_invalid"))
		return
	}

	// Отправляем сообщение об успешной смене языка
	s.ChannelMessageSend(m.ChannelID, localization.GetText("language_changed"))
}
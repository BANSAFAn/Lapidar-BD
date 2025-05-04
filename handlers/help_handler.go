package handlers

import (
	"fmt"
	"time"

	"discord-bot/localization"

	"github.com/bwmarrin/discordgo"
)

// HandleHelpCommand обрабатывает команду /help и отображает информацию о командах через вебхук
func HandleHelpCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Создаем вебхук в текущем канале
	webhook, err := s.WebhookCreate(m.ChannelID, "Lapidar Help", "")
	if err != nil {
		// Если не удалось создать вебхук, отправляем обычное сообщение
		s.ChannelMessageSend(m.ChannelID, localization.GetText("webhook_error"))

		// Отправляем справку через обычное сообщение
		sendHelpEmbed(s, m.ChannelID)
		return
	}

	// Создаем эмбед для справки
	embed := createHelpEmbed()

	// Отправляем сообщение через вебхук
	webhookParams := &discordgo.WebhookParams{
		Username:  "Lapidar Help",
		AvatarURL: s.State.User.AvatarURL(""),
		Embeds:    []*discordgo.MessageEmbed{embed},
	}

	// Отправляем сообщение через вебхук
	_, err = s.WebhookExecute(webhook.ID, webhook.Token, false, webhookParams)
	if err != nil {
		// Если не удалось отправить через вебхук, отправляем обычное сообщение
		s.ChannelMessageSend(m.ChannelID, localization.GetText("webhook_error"))
		sendHelpEmbed(s, m.ChannelID)
	}

	// Удаляем вебхук после использования
	s.WebhookDelete(webhook.ID)
}

// createHelpEmbed создает эмбед с информацией о командах
func createHelpEmbed() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: localization.GetText("help_title"),
		Color: 0x00BFFF,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  fmt.Sprintf("%sreport @пользователь причина", cfg.Prefix),
				Value: localization.GetText("report_command_desc"),
			},
			{
				Name:  fmt.Sprintf("%sban @пользователь причина [длительность]", cfg.Prefix),
				Value: localization.GetText("ban_command_desc"),
			},
			{
				Name:  fmt.Sprintf("%sai ваш запрос", cfg.Prefix),
				Value: localization.GetText("ai_command_desc"),
			},
			{
				Name:  fmt.Sprintf("%sgemini ваш запрос", cfg.Prefix),
				Value: localization.GetText("gemini_command_desc"),
			},
			{
				Name:  fmt.Sprintf("%sgrok ваш запрос", cfg.Prefix),
				Value: localization.GetText("grok_command_desc"),
			},
			{
				Name:  fmt.Sprintf("%schatgpt ваш запрос", cfg.Prefix),
				Value: localization.GetText("chatgpt_command_desc"),
			},
			{
				Name:  fmt.Sprintf("%sqwen ваш запрос", cfg.Prefix),
				Value: localization.GetText("qwen_command_desc"),
			},
			{
				Name:  fmt.Sprintf("%sclaude ваш запрос", cfg.Prefix),
				Value: localization.GetText("claude_command_desc"),
			},
			{
				Name:  fmt.Sprintf("%shelp", cfg.Prefix),
				Value: localization.GetText("help_command_desc"),
			},
			{
				Name:  fmt.Sprintf("%slanguage [ru|en|uk|de|zh]", cfg.Prefix),
				Value: localization.GetText("language_command_desc"),
			},
			{
				Name:  fmt.Sprintf("%splay URL-YouTube", cfg.Prefix),
				Value: localization.GetText("play_command_desc"),
			},
			{
				Name:  fmt.Sprintf("%sstop", cfg.Prefix),
				Value: localization.GetText("stop_command_desc"),
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Lapidar Bot",
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

// sendHelpEmbed отправляет эмбед с информацией о командах через обычное сообщение
func sendHelpEmbed(s *discordgo.Session, channelID string) {
	embed := createHelpEmbed()
	s.ChannelMessageSendEmbed(channelID, embed)
}

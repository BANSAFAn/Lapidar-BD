package handlers

import (
	"discord-bot/localization"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// HandleNicknameCommand обрабатывает команду для изменения никнейма пользователя
func HandleNicknameCommand(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) < 2 {
		if _, err := s.ChannelMessageSend(m.ChannelID, localization.GetText("nickname_usage", "!")); err != nil {
			fmt.Printf("Ошибка отправки сообщения: %v\n", err)
		}
		return
	}

	// Получаем ID пользователя из упоминания или первого аргумента
	userID := args[0]
	// Удаляем форматирование упоминания, если оно есть
	userID = extractUserID(userID)

	// Получаем новый никнейм из оставшихся аргументов
	newNickname := args[1]

	// Проверяем права бота на изменение никнеймов
	guild, err := s.State.Guild(m.GuildID)
	if err != nil {
		if _, err := s.ChannelMessageSend(m.ChannelID, localization.GetText("nickname_guild_error")); err != nil {
			fmt.Printf("Ошибка отправки сообщения: %v\n", err)
		}
		return
	}

	// Проверяем, имеет ли пользователь права на изменение никнеймов
	permissions, err := s.State.UserChannelPermissions(m.Author.ID, m.ChannelID)
	if err != nil || (permissions&discordgo.PermissionManageNicknames) == 0 {
		if _, err := s.ChannelMessageSend(m.ChannelID, localization.GetText("nickname_no_permission")); err != nil {
			fmt.Printf("Ошибка отправки сообщения: %v\n", err)
		}
		return
	}

	// Изменяем никнейм пользователя
	err = s.GuildMemberNickname(m.GuildID, userID, newNickname)
	if err != nil {
		if _, err := s.ChannelMessageSend(m.ChannelID, localization.GetText("nickname_error", err.Error())); err != nil {
			fmt.Printf("Ошибка отправки сообщения: %v\n", err)
		}
		return
	}

	if _, err := s.ChannelMessageSend(m.ChannelID, localization.GetText("nickname_success")); err != nil {
		fmt.Printf("Ошибка отправки сообщения: %v\n", err)
	}
}

// HandleDMCommand обрабатывает команду для отправки личного сообщения пользователю
func HandleDMCommand(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) < 2 {
		if _, err := s.ChannelMessageSend(m.ChannelID, localization.GetText("dm_usage", "!")); err != nil {
			fmt.Printf("Ошибка отправки сообщения: %v\n", err)
		}
		return
	}

	// Получаем ID пользователя из упоминания или первого аргумента
	userID := args[0]
	// Удаляем форматирование упоминания, если оно есть
	userID = extractUserID(userID)

	// Получаем сообщение из оставшихся аргументов
	message := args[1]

	// Создаем личный канал с пользователем
	channel, err := s.UserChannelCreate(userID)
	if err != nil {
		if _, err := s.ChannelMessageSend(m.ChannelID, localization.GetText("dm_channel_error", err.Error())); err != nil {
			fmt.Printf("Ошибка отправки сообщения: %v\n", err)
		}
		return
	}

	// Отправляем сообщение в личный канал
	_, err = s.ChannelMessageSend(channel.ID, message)
	if err != nil {
		if _, err := s.ChannelMessageSend(m.ChannelID, localization.GetText("dm_send_error", err.Error())); err != nil {
			fmt.Printf("Ошибка отправки сообщения: %v\n", err)
		}
		return
	}

	if _, err := s.ChannelMessageSend(m.ChannelID, localization.GetText("dm_success")); err != nil {
		fmt.Printf("Ошибка отправки сообщения: %v\n", err)
	}
}

// extractUserID извлекает ID пользователя из упоминания
func extractUserID(mention string) string {
	// Если это упоминание в формате <@!123456789>
	if len(mention) > 3 && mention[0:2] == "<@" && mention[len(mention)-1:] == ">" {
		// Удаляем <@! или <@ в начале и > в конце
		if mention[2:3] == "!" {
			return mention[3 : len(mention)-1]
		}
		return mention[2 : len(mention)-1]
	}
	// Если это просто ID
	return mention
}

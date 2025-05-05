package reports

import (
	"fmt"
	"sync"
	"time"

	"discord-bot/db"

	"github.com/bwmarrin/discordgo"
)

// ReportMessage содержит информацию о сообщении репорта
type ReportMessage struct {
	ReportID       int64  // ID репорта в базе данных
	MessageID      string // ID сообщения в Discord
	ReportedUserID string // ID пользователя, на которого пожаловались
	ReporterID     string // ID пользователя, который отправил жалобу
	Reason         string // Причина жалобы
}

var (
	// reportMessages хранит информацию о сообщениях репортов
	reportMessages = make(map[string]ReportMessage)
	reportMutex    sync.RWMutex
)

// CreateReport создает новый репорт и отправляет сообщение в канал модерации
func CreateReport(s *discordgo.Session, channelID, reportedUserID, reporterID, reason string) (int64, error) {
	// Добавляем репорт в базу данных
	reportID, err := db.AddReport(reportedUserID, reporterID, reason)
	if err != nil {
		return 0, fmt.Errorf("ошибка добавления репорта в базу данных: %w", err)
	}

	// Получаем информацию о пользователях
	reportedUser, err := s.User(reportedUserID)
	if err != nil {
		return 0, fmt.Errorf("ошибка получения информации о пользователе: %w", err)
	}

	reporter, err := s.User(reporterID)
	if err != nil {
		return 0, fmt.Errorf("ошибка получения информации о пользователе: %w", err)
	}

	// Создаем эмбед для репорта
	embed := &discordgo.MessageEmbed{
		Title: fmt.Sprintf("Репорт #%d", reportID),
		Color: 0xFFA500, // Оранжевый цвет для непроверенных репортов
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Пользователь",
				Value:  fmt.Sprintf("%s#%s (%s)", reportedUser.Username, reportedUser.Discriminator, reportedUserID),
				Inline: true,
			},
			{
				Name:   "Отправитель",
				Value:  fmt.Sprintf("%s#%s (%s)", reporter.Username, reporter.Discriminator, reporterID),
				Inline: true,
			},
			{
				Name:  "Причина",
				Value: reason,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Используйте реакции ✅ для подтверждения или ❌ для отклонения",
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	// Отправляем сообщение
	msg, err := s.ChannelMessageSendEmbed(channelID, embed)
	if err != nil {
		return 0, fmt.Errorf("ошибка отправки сообщения: %w", err)
	}

	// Добавляем реакции для модерации
	if err := s.MessageReactionAdd(channelID, msg.ID, "✅"); err != nil {
		fmt.Printf("Ошибка добавления реакции: %v\n", err)
	}
	if err := s.MessageReactionAdd(channelID, msg.ID, "❌"); err != nil {
		fmt.Printf("Ошибка добавления реакции: %v\n", err)
	}

	// Сохраняем информацию о сообщении
	reportMutex.Lock()
	reportMessages[msg.ID] = ReportMessage{
		ReportID:       reportID,
		MessageID:      msg.ID,
		ReportedUserID: reportedUserID,
		ReporterID:     reporterID,
		Reason:         reason,
	}
	reportMutex.Unlock()

	return reportID, nil
}

// GetReportMessage возвращает информацию о сообщении репорта по ID сообщения
func GetReportMessage(messageID string) (ReportMessage, bool) {
	reportMutex.RLock()
	defer reportMutex.RUnlock()

	report, exists := reportMessages[messageID]
	return report, exists
}

// RemoveReportMessage удаляет информацию о сообщении репорта
func RemoveReportMessage(messageID string) {
	reportMutex.Lock()
	delete(reportMessages, messageID)
	reportMutex.Unlock()
}

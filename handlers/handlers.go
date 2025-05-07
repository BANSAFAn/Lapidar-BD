package handlers

import (
	"fmt"
	"strings"
	"time"

	"discord-bot/config"
	"discord-bot/db"
	"discord-bot/localization"
	"discord-bot/reports"

	"github.com/bwmarrin/discordgo"
)

var cfg *config.Config

func init() {
	var err error
	cfg, err = config.Load()
	if err != nil {
		fmt.Println("Ошибка загрузки конфигурации в обработчиках:", err)
	}

	// Инициализируем систему локализации
	if err := localization.Initialize(); err != nil {
		fmt.Println("Ошибка инициализации системы локализации:", err)
	}
}

// MessageCreate обрабатывает входящие сообщения
// Поддерживает префикс '/' и имя бота Lapidar
func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Игнорируем сообщения от самого бота
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Проверяем, забанен ли пользователь
	banned, err := db.IsUserBanned(m.Author.ID)
	if err != nil {
		fmt.Println("Ошибка при проверке бана:", err)
	}

	if banned {
		// Удаляем сообщение от забаненного пользователя
		if err := s.ChannelMessageDelete(m.ChannelID, m.ID); err != nil {
			fmt.Printf("Ошибка при удалении сообщения: %v\n", err)
		}
		return
	}

	// Проверяем, начинается ли сообщение с префикса команды
	if !strings.HasPrefix(m.Content, cfg.Prefix) {
		return
	}

	// Разбиваем сообщение на команду и аргументы
	args := strings.Split(strings.TrimPrefix(m.Content, cfg.Prefix), " ")
	command := strings.ToLower(args[0])

	// Обработка команд
	switch command {
	case "report":
		handleReportCommand(s, m, args[1:])
	case "ban":
		handleBanCommand(s, m, args[1:])
	case "help":
		HandleHelpCommand(s, m)
	case "ai":
		HandleAICommand(s, m, args[1:])
	case "gemini":
		HandleAICommand(s, m, append([]string{"gemini"}, args[1:]...))
	case "grok":
		HandleAICommand(s, m, append([]string{"grok"}, args[1:]...))
	case "chatgpt":
		HandleAICommand(s, m, append([]string{"chatgpt"}, args[1:]...))
	case "qwen":
		HandleAICommand(s, m, append([]string{"qwen"}, args[1:]...))
	case "claude":
		HandleAICommand(s, m, append([]string{"claude"}, args[1:]...))
	case "language", "lang":
		HandleLanguageCommand(s, m, args[1:])
	case "play":
		HandlePlayCommand(s, m, args[1:])
	case "stop":
		HandleStopCommand(s, m)
	case "leave":
		HandleLeaveCommand(s, m)
	case "nickname", "nick":
		HandleNicknameCommand(s, m, args[1:])
	case "dm", "message":
		HandleDMCommand(s, m, args[1:])
	}
}

// ReactionAdd обрабатывает добавление реакций
func ReactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	// Игнорируем реакции от самого бота
	if r.UserID == s.State.User.ID {
		return
	}

	// Проверяем, является ли сообщение репортом
	reportMsg, exists := reports.GetReportMessage(r.MessageID)
	if !exists {
		return
	}

	// Проверяем, имеет ли пользователь права администратора или модератора
	member, err := s.GuildMember(r.GuildID, r.UserID)
	if err != nil {
		fmt.Println("Ошибка при получении информации о пользователе:", err)
		return
	}

	isAdmin := false
	for _, roleID := range member.Roles {
		if roleID == cfg.AdminRoleID || roleID == cfg.ModRoleID {
			isAdmin = true
			break
		}
	}

	if !isAdmin {
		return
	}

	// Обрабатываем реакции на репорт
	switch r.Emoji.Name {
	case "✅": // Подтверждение репорта
		handleReportConfirmation(s, r, reportMsg)
	case "❌": // Отклонение репорта
		handleReportRejection(s, r, reportMsg)
	}
}

// handleReportCommand обрабатывает команду репорта
func handleReportCommand(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) < 2 {
		s.ChannelMessageSend(m.ChannelID, localization.GetText("report_usage", cfg.Prefix))
		return
	}

	// Извлекаем ID пользователя из упоминания
	userID := strings.Trim(args[0], "<@!>")
	reason := strings.Join(args[1:], " ")

	// Создаем репорт
	reportID, err := reports.CreateReport(s, m.ChannelID, userID, m.Author.ID, reason)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, localization.GetText("report_error", err.Error()))
		return
	}

	s.ChannelMessageSend(m.ChannelID, localization.GetText("report_created", reportID))
}

// handleBanCommand обрабатывает команду бана
func handleBanCommand(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	// Проверяем права пользователя
	member, err := s.GuildMember(m.GuildID, m.Author.ID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Ошибка при получении информации о пользователе.")
		return
	}

	isAdmin := false
	for _, roleID := range member.Roles {
		if roleID == cfg.AdminRoleID || roleID == cfg.ModRoleID {
			isAdmin = true
			break
		}
	}

	if !isAdmin {
		s.ChannelMessageSend(m.ChannelID, localization.GetText("ban_no_permission"))
		return
	}

	if len(args) < 2 {
		s.ChannelMessageSend(m.ChannelID, localization.GetText("ban_usage", cfg.Prefix))
		return
	}

	// Извлекаем ID пользователя из упоминания
	userID := strings.Trim(args[0], "<@!>")

	// Определяем причину и длительность
	var duration *time.Duration
	reason := strings.Join(args[1:], " ")

	// Если указана длительность
	if len(args) > 2 && (strings.HasSuffix(args[len(args)-1], "d") ||
		strings.HasSuffix(args[len(args)-1], "h") ||
		strings.HasSuffix(args[len(args)-1], "m")) {

		durationStr := args[len(args)-1]
		reason = strings.Join(args[1:len(args)-1], " ")

		// Парсим длительность
		dur, err := time.ParseDuration(durationStr)
		if err == nil {
			duration = &dur
		}
	}

	// Баним пользователя
	err = db.AddBan(userID, reason, m.Author.ID, duration)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, localization.GetText("ban_error", err.Error()))
		return
	}

	// Отправляем сообщение о бане
	durationText := localization.GetText("ban_duration_forever")
	if duration != nil {
		durationText = localization.GetText("ban_duration_for", duration.String())
	}

	s.ChannelMessageSend(m.ChannelID, localization.GetText("ban_success", userID, durationText, reason))
}

// Используем новый обработчик команды help из help_handler.go

// handleReportConfirmation обрабатывает подтверждение репорта
func handleReportConfirmation(s *discordgo.Session, r *discordgo.MessageReactionAdd, reportMsg reports.ReportMessage) {
	// Подтверждаем репорт в базе данных
	err := db.ConfirmReport(reportMsg.ReportID, r.UserID)
	if err != nil {
		fmt.Println("Ошибка при подтверждении репорта:", err)
		return
	}

	// Получаем количество подтвержденных репортов
	count, err := db.GetReportCount(reportMsg.ReportedUserID)
	if err != nil {
		fmt.Println("Ошибка при получении количества репортов:", err)
		return
	}

	// Обновляем сообщение репорта
	embed := &discordgo.MessageEmbed{
		Title: fmt.Sprintf("Репорт #%d (Подтвержден)", reportMsg.ReportID),
		Color: 0x00FF00,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Пользователь",
				Value: fmt.Sprintf("<@%s>", reportMsg.ReportedUserID),
			},
			{
				Name:  "Причина",
				Value: reportMsg.Reason,
			},
			{
				Name:  "Отправитель",
				Value: fmt.Sprintf("<@%s>", reportMsg.ReporterID),
			},
			{
				Name:  "Подтвержден",
				Value: fmt.Sprintf("<@%s>", r.UserID),
			},
			{
				Name:  "Всего репортов",
				Value: fmt.Sprintf("%d/%d", count, cfg.ReportThreshold),
			},
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	if _, err := s.ChannelMessageEditEmbed(r.ChannelID, r.MessageID, embed); err != nil {
		fmt.Printf("Ошибка при редактировании сообщения: %v\n", err)
	}

	// Если достигнут порог репортов, баним пользователя
	if count >= cfg.ReportThreshold {
		reason := fmt.Sprintf("Автоматический бан по достижению порога репортов (%d)", cfg.ReportThreshold)
		var duration time.Duration = 7 * 24 * time.Hour // Бан на 7 дней

		err := db.AddBan(reportMsg.ReportedUserID, reason, s.State.User.ID, &duration)
		if err != nil {
			fmt.Println("Ошибка при автоматическом бане:", err)
			return
		}

		// Отправляем сообщение о бане
		s.ChannelMessageSend(r.ChannelID, fmt.Sprintf(
			"Пользователь <@%s> автоматически забанен на 7 дней по достижению порога репортов (%d).",
			reportMsg.ReportedUserID, cfg.ReportThreshold,
		))
	}
}

// handleReportRejection обрабатывает отклонение репорта
func handleReportRejection(s *discordgo.Session, r *discordgo.MessageReactionAdd, reportMsg reports.ReportMessage) {
	// Обновляем сообщение репорта
	embed := &discordgo.MessageEmbed{
		Title: fmt.Sprintf("Репорт #%d (Отклонен)", reportMsg.ReportID),
		Color: 0xFF0000,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Пользователь",
				Value: fmt.Sprintf("<@%s>", reportMsg.ReportedUserID),
			},
			{
				Name:  "Причина",
				Value: reportMsg.Reason,
			},
			{
				Name:  "Отправитель",
				Value: fmt.Sprintf("<@%s>", reportMsg.ReporterID),
			},
			{
				Name:  "Отклонен",
				Value: fmt.Sprintf("<@%s>", r.UserID),
			},
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	if _, err := s.ChannelMessageEditEmbed(r.ChannelID, r.MessageID, embed); err != nil {
		fmt.Printf("Ошибка при редактировании сообщения: %v\n", err)
	}

	// Удаляем репорт из кэша
	reports.RemoveReportMessage(r.MessageID)
}

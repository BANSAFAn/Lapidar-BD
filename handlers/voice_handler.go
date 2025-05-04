package handlers

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"discord-bot/localization"

	"github.com/bwmarrin/discordgo"
	"github.com/kkdai/youtube/v2"
)

// Структура для хранения информации о голосовых соединениях
type VoiceInstance struct {
	connection *discordgo.VoiceConnection
	guildID    string
	channelID  string
	stopped    bool
	mutex      sync.Mutex
}

// Карта активных голосовых соединений
var voiceInstances = make(map[string]*VoiceInstance)

// Мьютекс для безопасного доступа к карте голосовых соединений
var voiceMutex sync.Mutex

// HandlePlayCommand обрабатывает команду воспроизведения аудио с YouTube
func HandlePlayCommand(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	// Проверяем, что пользователь указал URL
	if len(args) < 1 {
		s.ChannelMessageSend(m.ChannelID, localization.GetText("play_usage", cfg.Prefix))
		return
	}

	// Получаем URL YouTube
	url := args[0]

	// Проверяем, что URL валидный
	if !isValidYouTubeURL(url) {
		s.ChannelMessageSend(m.ChannelID, localization.GetText("play_invalid_url"))
		return
	}

	// Находим голосовой канал, в котором находится пользователь
	voiceChannelID := findUserVoiceChannel(s, m.GuildID, m.Author.ID)
	if voiceChannelID == "" {
		s.ChannelMessageSend(m.ChannelID, localization.GetText("play_not_in_voice"))
		return
	}

	// Отправляем сообщение о подключении к голосовому каналу
	s.ChannelMessageSend(m.ChannelID, localization.GetText("play_joining"))

	// Подключаемся к голосовому каналу
	vc, err := joinVoiceChannel(s, m.GuildID, voiceChannelID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, localization.GetText("play_error", err.Error()))
		return
	}

	// Получаем информацию о видео
	videoTitle, err := playYouTubeAudio(s, vc, url, m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, localization.GetText("play_error", err.Error()))
		return
	}

	// Отправляем сообщение о том, что сейчас играет
	s.ChannelMessageSend(m.ChannelID, localization.GetText("play_now_playing", videoTitle))
}

// findUserVoiceChannel находит голосовой канал, в котором находится пользователь
func findUserVoiceChannel(s *discordgo.Session, guildID, userID string) string {
	// Получаем список голосовых состояний для гильдии
	guild, err := s.State.Guild(guildID)
	if err != nil {
		return ""
	}

	// Ищем пользователя в голосовых каналах
	for _, vs := range guild.VoiceStates {
		if vs.UserID == userID {
			return vs.ChannelID
		}
	}

	return ""
}

// joinVoiceChannel подключается к голосовому каналу
func joinVoiceChannel(s *discordgo.Session, guildID, channelID string) (*VoiceInstance, error) {
	voiceMutex.Lock()
	defer voiceMutex.Unlock()

	// Проверяем, есть ли уже активное соединение для этой гильдии
	if vi, exists := voiceInstances[guildID]; exists {
		// Если бот уже в этом канале, возвращаем существующее соединение
		if vi.channelID == channelID {
			return vi, nil
		}

		// Если бот в другом канале, отключаемся от него
		vi.connection.Disconnect()
		delete(voiceInstances, guildID)
	}

	// Подключаемся к голосовому каналу
	vc, err := s.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		return nil, err
	}

	// Создаем новый экземпляр голосового соединения
	vi := &VoiceInstance{
		connection: vc,
		guildID:    guildID,
		channelID:  channelID,
		stopped:    false,
	}

	// Сохраняем соединение в карте
	voiceInstances[guildID] = vi

	return vi, nil
}

// playYouTubeAudio воспроизводит аудио с YouTube
func playYouTubeAudio(s *discordgo.Session, vi *VoiceInstance, url, guildID string) (string, error) {
	// Создаем клиент YouTube
	client := youtube.Client{}

	// Получаем информацию о видео
	video, err := client.GetVideo(url)
	if err != nil {
		return "", fmt.Errorf("ошибка получения информации о видео: %w", err)
	}

	// Получаем форматы только с аудио
	formats := video.Formats.WithAudioChannels()
	if len(formats) == 0 {
		return "", fmt.Errorf("не найдены аудио форматы для видео")
	}

	// Выбираем формат с лучшим качеством аудио
	format := formats.FindByQuality("tiny")
	if format.Empty() {
		format = formats[0]
	}

	// Получаем URL для скачивания
	downloadURL, err := client.GetStreamURL(video, &format)
	if err != nil {
		return "", fmt.Errorf("ошибка получения URL для скачивания: %w", err)
	}

	// Скачиваем аудио
	resp, err := http.Get(downloadURL)
	if err != nil {
		return "", fmt.Errorf("ошибка скачивания аудио: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем, что соединение все еще активно
	voiceMutex.Lock()
	vi, exists := voiceInstances[guildID]
	voiceMutex.Unlock()

	if !exists || vi.stopped {
		return "", fmt.Errorf("голосовое соединение было закрыто")
	}

	// Начинаем говорить (это необходимо для отправки аудио)
	vi.connection.Speaking(true)
	defer vi.connection.Speaking(false)

	// Создаем буфер для аудио данных
	buffer := make([]byte, 16384)
	for {
		// Проверяем, не было ли остановлено воспроизведение
		vi.mutex.Lock()
		if vi.stopped {
			vi.mutex.Unlock()
			break
		}
		vi.mutex.Unlock()

		// Читаем данные из ответа
		n, err := resp.Body.Read(buffer)
		if err != nil && err != io.EOF {
			return video.Title, fmt.Errorf("ошибка чтения аудио данных: %w", err)
		}

		if n == 0 || err == io.EOF {
			break
		}

		// Отправляем аудио данные в голосовое соединение
		vi.connection.OpusSend <- buffer[:n]
	}

	return video.Title, nil
}

// isValidYouTubeURL проверяет, является ли URL валидным URL YouTube
func isValidYouTubeURL(url string) bool {
	// Регулярное выражение для проверки URL YouTube
	youtubeRegex := regexp.MustCompile(`^(https?://)?(www\.)?(youtube\.com|youtu\.?be)/.+$`)
	return youtubeRegex.MatchString(url)
}

// StopVoice останавливает воспроизведение и отключается от голосового канала
func StopVoice(s *discordgo.Session, guildID string) error {
	voiceMutex.Lock()
	defer voiceMutex.Unlock()

	// Проверяем, есть ли активное соединение для этой гильдии
	vi, exists := voiceInstances[guildID]
	if !exists {
		return fmt.Errorf("нет активного голосового соединения")
	}

	// Останавливаем воспроизведение
	vi.mutex.Lock()
	vi.stopped = true
	vi.mutex.Unlock()

	// Отключаемся от голосового канала
	err := vi.connection.Disconnect()

	// Удаляем соединение из карты
	delete(voiceInstances, guildID)

	return err
}

// HandleStopCommand обрабатывает команду остановки воспроизведения
func HandleStopCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Останавливаем воспроизведение
	err := StopVoice(s, m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, localization.GetText("stop_error", err.Error()))
		return
	}

	// Отправляем сообщение об успешной остановке
	s.ChannelMessageSend(m.ChannelID, localization.GetText("stop_success"))
}
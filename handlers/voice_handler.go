package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"sync"

	"discord-bot/localization"

	"github.com/bwmarrin/discordgo"
	"github.com/kkdai/youtube/v2"
)

type VoiceInstance struct {
	connection *discordgo.VoiceConnection
	guildID    string
	channelID  string
	stopped    bool
	mutex      sync.Mutex
}

var voiceInstances = make(map[string]*VoiceInstance)
var voiceMutex sync.Mutex

func HandlePlayCommand(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) < 1 {
		if _, err := s.ChannelMessageSend(m.ChannelID, localization.GetText("play_usage", "!")); err != nil {
			fmt.Printf("Ошибка отправки сообщения: %v\n", err)
		}
		return
	}

	url := args[0]

	if !isValidYouTubeURL(url) {
		if _, err := s.ChannelMessageSend(m.ChannelID, localization.GetText("play_invalid_url")); err != nil {
			fmt.Printf("Ошибка отправки сообщения: %v\n", err)
		}
		return
	}

	voiceChannelID := findUserVoiceChannel(s, m.GuildID, m.Author.ID)
	if voiceChannelID == "" {
		if _, err := s.ChannelMessageSend(m.ChannelID, localization.GetText("play_not_in_voice")); err != nil {
			fmt.Printf("Ошибка отправки сообщения: %v\n", err)
		}
		return
	}

	if _, err := s.ChannelMessageSend(m.ChannelID, localization.GetText("play_joining")); err != nil {
		fmt.Printf("Ошибка отправки сообщения: %v\n", err)
	}

	vc, err := joinVoiceChannel(s, m.GuildID, voiceChannelID)
	if err != nil {
		if _, err2 := s.ChannelMessageSend(m.ChannelID, localization.GetText("play_error", err.Error())); err2 != nil {
			fmt.Printf("Ошибка отправки сообщения: %v\n", err2)
		}
		return
	}

	videoTitle, err := playYouTubeAudio(s, vc, url, m.GuildID)
	if err != nil {
		if _, err2 := s.ChannelMessageSend(m.ChannelID, localization.GetText("play_error", err.Error())); err2 != nil {
			fmt.Printf("Ошибка отправки сообщения: %v\n", err2)
		}
		return
	}

	if _, err := s.ChannelMessageSend(m.ChannelID, localization.GetText("play_now_playing", videoTitle)); err != nil {
		fmt.Printf("Ошибка отправки сообщения: %v\n", err)
	}
}

func findUserVoiceChannel(s *discordgo.Session, guildID, userID string) string {
	guild, err := s.State.Guild(guildID)
	if err != nil {
		return ""
	}

	for _, vs := range guild.VoiceStates {
		if vs.UserID == userID {
			return vs.ChannelID
		}
	}

	return ""
}

func joinVoiceChannel(s *discordgo.Session, guildID, channelID string) (*VoiceInstance, error) {
	voiceMutex.Lock()
	defer voiceMutex.Unlock()

	if vi, exists := voiceInstances[guildID]; exists {
		if vi.channelID == channelID {
			return vi, nil
		}

		if err := vi.connection.Disconnect(); err != nil {
			return nil, fmt.Errorf("ошибка при отключении от голосового канала: %w", err)
		}
		delete(voiceInstances, guildID)
	}

	vc, err := s.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		return nil, err
	}

	vi := &VoiceInstance{
		connection: vc,
		guildID:    guildID,
		channelID:  channelID,
		stopped:    false,
	}

	voiceInstances[guildID] = vi

	return vi, nil
}

func playYouTubeAudio(s *discordgo.Session, vi *VoiceInstance, url, guildID string) (string, error) {
	client := youtube.Client{}

	video, err := client.GetVideo(url)
	if err != nil {
		return "", fmt.Errorf("ошибка получения информации о видео: %w", err)
	}

	formats := video.Formats.WithAudioChannels()
	if len(formats) == 0 {
		return "", fmt.Errorf("не найдены аудио форматы для видео")
	}

	var format youtube.Format
	formatPtr := formats.FindByQuality("tiny")
	if formatPtr == nil {
		format = formats[0]
	} else {
		format = *formatPtr
	}

	resp, _, err := client.GetStream(video, &format)
	if err != nil {
		return "", fmt.Errorf("ошибка получения потока: %w", err)
	}
	defer resp.Close()

	vi.mutex.Lock()
	vi.stopped = false
	vi.mutex.Unlock()

	vi.connection.Speaking(true)
	defer vi.connection.Speaking(false)

	buffer := make([]byte, 16384)
	for {
		vi.mutex.Lock()
		if vi.stopped {
			vi.mutex.Unlock()
			break
		}
		vi.mutex.Unlock()

		n, err := resp.Read(buffer)
		if err != nil && err != io.EOF {
			return video.Title, fmt.Errorf("ошибка чтения потока: %w", err)
		}

		if n == 0 {
			break
		}

		vi.connection.OpusSend <- buffer[:n]
	}

	return video.Title, nil
}

func isValidYouTubeURL(url string) bool {
	youtubeRegex := regexp.MustCompile(`^(https?://)?(www\.)?(youtube\.com|youtu\.?be)/.+$`)
	return youtubeRegex.MatchString(url)
}

func HandleStopCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	voiceMutex.Lock()
	vi, exists := voiceInstances[m.GuildID]
	voiceMutex.Unlock()

	if !exists {
		if _, err := s.ChannelMessageSend(m.ChannelID, localization.GetText("stop_not_playing")); err != nil {
			fmt.Printf("Ошибка отправки сообщения: %v\n", err)
		}
		return
	}

	vi.mutex.Lock()
	vi.stopped = true
	vi.mutex.Unlock()

	if _, err := s.ChannelMessageSend(m.ChannelID, localization.GetText("stop_success")); err != nil {
		fmt.Printf("Ошибка отправки сообщения: %v\n", err)
	}
}

func HandleLeaveCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	voiceMutex.Lock()
	defer voiceMutex.Unlock()

	vi, exists := voiceInstances[m.GuildID]
	if !exists {
		if _, err := s.ChannelMessageSend(m.ChannelID, localization.GetText("leave_not_in_voice")); err != nil {
			fmt.Printf("Ошибка отправки сообщения: %v\n", err)
		}
		return
	}

	vi.mutex.Lock()
	vi.stopped = true
	vi.mutex.Unlock()

	if err := vi.connection.Disconnect(); err != nil {
		if _, err2 := s.ChannelMessageSend(m.ChannelID, localization.GetText("leave_error", err.Error())); err2 != nil {
			fmt.Printf("Ошибка отправки сообщения об ошибке: %v\n", err2)
		}
		return
	}
	delete(voiceInstances, m.GuildID)

	if _, err := s.ChannelMessageSend(m.ChannelID, localization.GetText("leave_success")); err != nil {
		fmt.Printf("Ошибка отправки сообщения: %v\n", err)
	}
}

func DownloadYouTubeAudio(url string) (string, error) {
	client := youtube.Client{}

	video, err := client.GetVideo(url)
	if err != nil {
		return "", fmt.Errorf("ошибка получения информации о видео: %w", err)
	}

	formats := video.Formats.WithAudioChannels()
	if len(formats) == 0 {
		return "", fmt.Errorf("не найдены аудио форматы для видео")
	}

	var format youtube.Format
	formatPtr := formats.FindByQuality("tiny")
	if formatPtr == nil {
		format = formats[0]
	} else {
		format = *formatPtr
	}

	// Получаем URL для скачивания напрямую
	url, err = client.GetStreamURL(video, &format)
	if err != nil {
		return "", fmt.Errorf("ошибка получения URL потока: %w", err)
	}

	// Создаем директорию для аудио файлов, если она не существует
	if err := os.MkdirAll("data/audio", 0755); err != nil {
		return "", fmt.Errorf("ошибка создания директории: %w", err)
	}

	filePath := fmt.Sprintf("data/audio/%s.mp3", video.ID)

	// Получаем данные по URL
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("ошибка загрузки файла: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ошибка при загрузке файла: статус %d", resp.StatusCode)
	}

	// Создаем файл для сохранения
	out, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("ошибка создания файла: %w", err)
	}
	defer out.Close()

	// Копируем данные из ответа в файл
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", fmt.Errorf("ошибка при записи файла: %w", err)
	}

	return filePath, nil
}

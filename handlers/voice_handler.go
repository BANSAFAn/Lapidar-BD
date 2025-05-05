package handlers

import (
	"fmt"
	"io"
	"net/http"
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
		s.ChannelMessageSend(m.ChannelID, localization.GetText("play_usage", "!"))
		return
	}

	url := args[0]

	if !isValidYouTubeURL(url) {
		s.ChannelMessageSend(m.ChannelID, localization.GetText("play_invalid_url"))
		return
	}

	voiceChannelID := findUserVoiceChannel(s, m.GuildID, m.Author.ID)
	if voiceChannelID == "" {
		s.ChannelMessageSend(m.ChannelID, localization.GetText("play_not_in_voice"))
		return
	}

	s.ChannelMessageSend(m.ChannelID, localization.GetText("play_joining"))

	vc, err := joinVoiceChannel(s, m.GuildID, voiceChannelID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, localization.GetText("play_error", err.Error()))
		return
	}

	videoTitle, err := playYouTubeAudio(s, vc, url, m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, localization.GetText("play_error", err.Error()))
		return
	}

	s.ChannelMessageSend(m.ChannelID, localization.GetText("play_now_playing", videoTitle))
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

		vi.connection.Disconnect()
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
		s.ChannelMessageSend(m.ChannelID, localization.GetText("stop_not_playing"))
		return
	}

	vi.mutex.Lock()
	vi.stopped = true
	vi.mutex.Unlock()

	s.ChannelMessageSend(m.ChannelID, localization.GetText("stop_success"))
}

func HandleLeaveCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	voiceMutex.Lock()
	defer voiceMutex.Unlock()

	vi, exists := voiceInstances[m.GuildID]
	if !exists {
		s.ChannelMessageSend(m.ChannelID, localization.GetText("leave_not_in_voice"))
		return
	}

	vi.mutex.Lock()
	vi.stopped = true
	vi.mutex.Unlock()

	vi.connection.Disconnect()
	delete(voiceInstances, m.GuildID)

	s.ChannelMessageSend(m.ChannelID, localization.GetText("leave_success"))
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

	filePath := fmt.Sprintf("data/audio/%s.mp3", video.ID)

	file, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("ошибка загрузки файла: %w", err)
	}
	defer file.Body.Close()

	return filePath, nil
}

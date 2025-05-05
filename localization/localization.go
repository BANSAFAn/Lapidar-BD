package localization

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"discord-bot/config"
)

// Поддерживаемые языки
const (
	Russian     = "ru"
	English     = "en"
	Ukrainian   = "uk"
	German      = "de"
	ChineseSimp = "zh"
)

// Локализованные строки для каждого языка
var translations map[string]map[string]string

// Текущий язык бота
var currentLanguage string

// Initialize инициализирует систему локализации
func Initialize() error {
	// Загружаем конфигурацию
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("ошибка загрузки конфигурации: %w", err)
	}

	// Устанавливаем язык по умолчанию из конфигурации или используем русский
	currentLanguage = cfg.DefaultLanguage
	if currentLanguage == "" {
		currentLanguage = Russian
	}

	// Создаем директорию для локализаций, если она не существует
	localizationDir := "localization/translations"
	if err := os.MkdirAll(localizationDir, 0755); err != nil {
		return fmt.Errorf("ошибка создания директории локализации: %w", err)
	}

	// Инициализируем карту переводов
	translations = make(map[string]map[string]string)

	// Загружаем все файлы локализации
	return loadTranslations(localizationDir)
}

// loadTranslations загружает все файлы локализации из указанной директории
func loadTranslations(dir string) error {
	// Проверяем наличие файлов локализации
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		// Если директория не существует или пуста, создаем файлы локализации по умолчанию
		if os.IsNotExist(err) || len(files) == 0 {
			return createDefaultTranslations(dir)
		}
		return fmt.Errorf("ошибка чтения директории локализации: %w", err)
	}

	// Загружаем каждый файл локализации
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		// Извлекаем код языка из имени файла (например, "ru.json" -> "ru")
		langCode := strings.TrimSuffix(file.Name(), ".json")

		// Загружаем файл локализации
		filePath := filepath.Join(dir, file.Name())
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("ошибка чтения файла локализации %s: %w", file.Name(), err)
		}

		// Декодируем JSON
		var langTranslations map[string]string
		if err := json.Unmarshal(data, &langTranslations); err != nil {
			return fmt.Errorf("ошибка декодирования файла локализации %s: %w", file.Name(), err)
		}

		// Сохраняем переводы для этого языка
		translations[langCode] = langTranslations
	}

	// Проверяем, что все необходимые языки загружены
	requiredLanguages := []string{Russian, English, Ukrainian, German, ChineseSimp}
	for _, lang := range requiredLanguages {
		if _, exists := translations[lang]; !exists {
			// Если какого-то языка нет, создаем файл локализации по умолчанию для него
			if err := createTranslationFile(dir, lang); err != nil {
				return fmt.Errorf("ошибка создания файла локализации для %s: %w", lang, err)
			}
		}
	}

	return nil
}

// createDefaultTranslations создает файлы локализации по умолчанию для всех поддерживаемых языков
func createDefaultTranslations(dir string) error {
	languages := []string{Russian, English, Ukrainian, German, ChineseSimp}
	for _, lang := range languages {
		if err := createTranslationFile(dir, lang); err != nil {
			return err
		}
	}
	return nil
}

// createTranslationFile создает файл локализации для указанного языка
func createTranslationFile(dir string, lang string) error {
	// Получаем переводы по умолчанию для указанного языка
	defaultTranslations := getDefaultTranslations(lang)

	// Сериализуем в JSON
	data, err := json.MarshalIndent(defaultTranslations, "", "  ")
	if err != nil {
		return fmt.Errorf("ошибка сериализации переводов для %s: %w", lang, err)
	}

	// Создаем файл
	filePath := filepath.Join(dir, lang+".json")
	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("ошибка записи файла локализации для %s: %w", lang, err)
	}

	// Сохраняем переводы в памяти
	translations[lang] = defaultTranslations

	return nil
}

// getDefaultTranslations возвращает переводы по умолчанию для указанного языка
func getDefaultTranslations(lang string) map[string]string {
	switch lang {
	case Russian:
		return map[string]string{
			"help_title":            "Справка по командам",
			"report_command_desc":   "Отправить жалобу на пользователя",
			"ban_command_desc":      "Забанить пользователя (только для администраторов)",
			"ai_command_desc":       "Задать вопрос искусственному интеллекту Gemini",
			"help_command_desc":     "Показать эту справку",
			"language_command_desc": "Изменить язык бота",
			"play_command_desc":     "Воспроизвести аудио с YouTube в голосовом канале",
			"report_usage":          "Использование: %sreport @пользователь причина",
			"ban_usage":             "Использование: %sban @пользователь причина [длительность]",
			"ai_usage":              "Использование: %sai ваш запрос",
			"language_usage":        "Использование: %slanguage [ru|en|uk|de|zh]",
			"play_usage":            "Использование: %splay URL-YouTube",
			"report_created":        "Репорт #%d создан и отправлен на рассмотрение администрации.",
			"report_error":          "Ошибка при создании репорта: %s",
			"ban_no_permission":     "У вас нет прав для использования этой команды.",
			"ban_error":             "Ошибка при добавлении бана: %s",
			"ban_success":           "Пользователь <@%s> забанен %s. Причина: %s",
			"ban_duration_forever":  "навсегда",
			"ban_duration_for":      "на %s",
			"ai_processing":         "Обрабатываю запрос, пожалуйста, подождите...",
			"ai_error":              "Ошибка при обработке запроса: %s",
			"language_changed":      "Язык бота изменен на русский.",
			"language_invalid":      "Неверный код языка. Доступные языки: ru (русский), en (английский), uk (украинский), de (немецкий), zh (китайский).",
			"play_not_in_voice":     "Вы должны находиться в голосовом канале.",
			"play_joining":          "Присоединяюсь к голосовому каналу...",
			"play_error":            "Ошибка при воспроизведении: %s",
			"play_now_playing":      "Сейчас играет: %s",
			"play_invalid_url":      "Неверный URL YouTube.",
		}
	case English:
		return map[string]string{
			"help_title":            "Command Help",
			"report_command_desc":   "Report a user",
			"ban_command_desc":      "Ban a user (administrators only)",
			"ai_command_desc":       "Ask a question to Gemini AI",
			"help_command_desc":     "Show this help",
			"language_command_desc": "Change bot language",
			"play_command_desc":     "Play audio from YouTube in a voice channel",
			"report_usage":          "Usage: %sreport @user reason",
			"ban_usage":             "Usage: %sban @user reason [duration]",
			"ai_usage":              "Usage: %sai your query",
			"language_usage":        "Usage: %slanguage [ru|en|uk|de|zh]",
			"play_usage":            "Usage: %splay YouTube-URL",
			"report_created":        "Report #%d created and sent to administrators.",
			"report_error":          "Error creating report: %s",
			"ban_no_permission":     "You don't have permission to use this command.",
			"ban_error":             "Error adding ban: %s",
			"ban_success":           "User <@%s> has been banned %s. Reason: %s",
			"ban_duration_forever":  "forever",
			"ban_duration_for":      "for %s",
			"ai_processing":         "Processing request, please wait...",
			"ai_error":              "Error processing request: %s",
			"language_changed":      "Bot language changed to English.",
			"language_invalid":      "Invalid language code. Available languages: ru (Russian), en (English), uk (Ukrainian), de (German), zh (Chinese).",
			"play_not_in_voice":     "You must be in a voice channel.",
			"play_joining":          "Joining voice channel...",
			"play_error":            "Error playing audio: %s",
			"play_now_playing":      "Now playing: %s",
			"play_invalid_url":      "Invalid YouTube URL.",
		}
	case Ukrainian:
		return map[string]string{
			"help_title":            "Довідка по командам",
			"report_command_desc":   "Відправити скаргу на користувача",
			"ban_command_desc":      "Заблокувати користувача (тільки для адміністраторів)",
			"ai_command_desc":       "Задати питання штучному інтелекту Gemini",
			"help_command_desc":     "Показати цю довідку",
			"language_command_desc": "Змінити мову бота",
			"play_command_desc":     "Відтворити аудіо з YouTube у голосовому каналі",
			"report_usage":          "Використання: %sreport @користувач причина",
			"ban_usage":             "Використання: %sban @користувач причина [тривалість]",
			"ai_usage":              "Використання: %sai ваш запит",
			"language_usage":        "Використання: %slanguage [ru|en|uk|de|zh]",
			"play_usage":            "Використання: %splay URL-YouTube",
			"report_created":        "Скарга #%d створена і відправлена на розгляд адміністрації.",
			"report_error":          "Помилка при створенні скарги: %s",
			"ban_no_permission":     "У вас немає прав для використання цієї команди.",
			"ban_error":             "Помилка при додаванні блокування: %s",
			"ban_success":           "Користувач <@%s> заблокований %s. Причина: %s",
			"ban_duration_forever":  "назавжди",
			"ban_duration_for":      "на %s",
			"ai_processing":         "Обробляю запит, будь ласка, зачекайте...",
			"ai_error":              "Помилка при обробці запиту: %s",
			"language_changed":      "Мову бота змінено на українську.",
			"language_invalid":      "Невірний код мови. Доступні мови: ru (російська), en (англійська), uk (українська), de (німецька), zh (китайська).",
			"play_not_in_voice":     "Ви повинні знаходитися в голосовому каналі.",
			"play_joining":          "Приєднуюсь до голосового каналу...",
			"play_error":            "Помилка при відтворенні: %s",
			"play_now_playing":      "Зараз грає: %s",
			"play_invalid_url":      "Невірний URL YouTube.",
		}
	case German:
		return map[string]string{
			"help_title":            "Befehlshilfe",
			"report_command_desc":   "Einen Benutzer melden",
			"ban_command_desc":      "Einen Benutzer sperren (nur für Administratoren)",
			"ai_command_desc":       "Stelle eine Frage an die Gemini KI",
			"help_command_desc":     "Diese Hilfe anzeigen",
			"language_command_desc": "Bot-Sprache ändern",
			"play_command_desc":     "Audio von YouTube in einem Sprachkanal abspielen",
			"report_usage":          "Verwendung: %sreport @Benutzer Grund",
			"ban_usage":             "Verwendung: %sban @Benutzer Grund [Dauer]",
			"ai_usage":              "Verwendung: %sai deine Anfrage",
			"language_usage":        "Verwendung: %slanguage [ru|en|uk|de|zh]",
			"play_usage":            "Verwendung: %splay YouTube-URL",
			"report_created":        "Meldung #%d erstellt und an Administratoren gesendet.",
			"report_error":          "Fehler beim Erstellen der Meldung: %s",
			"ban_no_permission":     "Du hast keine Berechtigung, diesen Befehl zu verwenden.",
			"ban_error":             "Fehler beim Hinzufügen der Sperre: %s",
			"ban_success":           "Benutzer <@%s> wurde %s gesperrt. Grund: %s",
			"ban_duration_forever":  "für immer",
			"ban_duration_for":      "für %s",
			"ai_processing":         "Anfrage wird verarbeitet, bitte warten...",
			"ai_error":              "Fehler bei der Verarbeitung der Anfrage: %s",
			"language_changed":      "Bot-Sprache wurde auf Deutsch geändert.",
			"language_invalid":      "Ungültiger Sprachcode. Verfügbare Sprachen: ru (Russisch), en (Englisch), uk (Ukrainisch), de (Deutsch), zh (Chinesisch).",
			"play_not_in_voice":     "Du musst dich in einem Sprachkanal befinden.",
			"play_joining":          "Trete dem Sprachkanal bei...",
			"play_error":            "Fehler beim Abspielen: %s",
			"play_now_playing":      "Spielt jetzt: %s",
			"play_invalid_url":      "Ungültige YouTube-URL.",
		}
	case ChineseSimp:
		return map[string]string{
			"help_title":            "命令帮助",
			"report_command_desc":   "举报用户",
			"ban_command_desc":      "封禁用户（仅限管理员）",
			"ai_command_desc":       "向Gemini人工智能提问",
			"help_command_desc":     "显示此帮助",
			"language_command_desc": "更改机器人语言",
			"play_command_desc":     "在语音频道中播放YouTube音频",
			"report_usage":          "用法: %sreport @用户 原因",
			"ban_usage":             "用法: %sban @用户 原因 [时长]",
			"ai_usage":              "用法: %sai 您的问题",
			"language_usage":        "用法: %slanguage [ru|en|uk|de|zh]",
			"play_usage":            "用法: %splay YouTube-URL",
			"report_created":        "举报 #%d 已创建并发送给管理员。",
			"report_error":          "创建举报时出错: %s",
			"ban_no_permission":     "您没有使用此命令的权限。",
			"ban_error":             "添加封禁时出错: %s",
			"ban_success":           "用户 <@%s> 已被封禁 %s。原因: %s",
			"ban_duration_forever":  "永久",
			"ban_duration_for":      "%s",
			"ai_processing":         "正在处理请求，请稍候...",
			"ai_error":              "处理请求时出错: %s",
			"language_changed":      "机器人语言已更改为简体中文。",
			"language_invalid":      "无效的语言代码。可用语言: ru (俄语), en (英语), uk (乌克兰语), de (德语), zh (中文)。",
			"play_not_in_voice":     "您必须在语音频道中。",
			"play_joining":          "正在加入语音频道...",
			"play_error":            "播放时出错: %s",
			"play_now_playing":      "正在播放: %s",
			"play_invalid_url":      "无效的YouTube URL。",
		}
	default:
		// По умолчанию возвращаем английские переводы
		return getDefaultTranslations(English)
	}
}

// GetText возвращает локализованный текст для указанного ключа
func GetText(key string, args ...interface{}) string {
	// Получаем переводы для текущего языка
	langTranslations, exists := translations[currentLanguage]
	if !exists {
		// Если переводы для текущего языка не найдены, используем английский
		langTranslations = translations[English]
		if langTranslations == nil {
			// Если и английские переводы не найдены, возвращаем ключ
			return key
		}
	}

	// Получаем перевод для указанного ключа
	text, exists := langTranslations[key]
	if !exists {
		// Если перевод не найден, пробуем найти его в английских переводах
		if currentLanguage != English {
			enTranslations := translations[English]
			if enTranslations != nil {
				text, exists = enTranslations[key]
				if !exists {
					// Если и в английских переводах не найден, возвращаем ключ
					return key
				}
			} else {
				return key
			}
		} else {
			return key
		}
	}

	// Если есть аргументы, форматируем строку
	if len(args) > 0 {
		return fmt.Sprintf(text, args...)
	}

	return text
}

// SetLanguage устанавливает текущий язык бота
func SetLanguage(lang string) bool {
	// Проверяем, что указанный язык поддерживается
	if _, exists := translations[lang]; !exists {
		return false
	}

	// Устанавливаем текущий язык
	currentLanguage = lang

	// Обновляем конфигурацию
	cfg, err := config.Load()
	if err == nil {
		cfg.DefaultLanguage = lang
		if err := config.SaveConfig(cfg); err != nil {
			fmt.Printf("Ошибка сохранения конфигурации: %v\n", err)
		}
	}

	return true
}

// GetCurrentLanguage возвращает текущий язык бота
func GetCurrentLanguage() string {
	return currentLanguage
}

// GetAvailableLanguages возвращает список доступных языков
func GetAvailableLanguages() []string {
	return []string{Russian, English, Ukrainian, German, ChineseSimp}
}

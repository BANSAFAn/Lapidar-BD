# Lapidar - Multifunctional Discord Bot in Go

<div align="center">

<h1 style="border: none; font-size: 2.5em;">LAPIDAR</h1>

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go&logoColor=white)](https://golang.org/)
[![Discord](https://img.shields.io/badge/Discord-7289DA?logo=discord&logoColor=white)](https://discord.gg/5YJ2uR5tDq)
[![SQLite](https://img.shields.io/badge/SQLite-003B57?logo=sqlite&logoColor=white)](https://www.sqlite.org/)
[![CI/CD Pipeline](https://github.com/BANSAFAn/Lapidar-BD/actions/workflows/ci-cd.yml/badge.svg)](https://github.com/BANSAFAn/Lapidar-BD/actions/workflows/ci-cd.yml)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/BANSAFAn/Lapidar-BD)](https://github.com/BANSAFAn/Lapidar-BD/releases)

</div>

<div align="center">

**Language Versions:** 
[English](README/README.en.md) | 
[Русский](README/README.ru.md) | 
[Українська](README/README.uk.md) | 
[Deutsch](README/README.de.md) | 
[中文](README/README.zh.md)

</div>

Lapidar is a powerful multifunctional Discord bot written in Go. It offers a wide range of server management capabilities, including a moderation system, Gemini AI integration, multilingual support, and YouTube audio playback in voice channels.

## 🌟 Key Features

- **Advanced Report System**
  - Submit reports on users with administrator confirmation
  - Automatic ban when reaching the report threshold
  - Detailed report history for each user

- **Moderation Tools**
  - Commands to ban users with reason and duration
  - Flexible role system for administrators and moderators
  - Logging of all moderation actions

- **AI Integration**
  - Multiple AI models (Gemini, Grok, ChatGPT, Qwen, Claude)
  - Ask questions to artificial intelligence directly in Discord
  - Get informative and accurate answers to your queries
  - Customizable parameters for response generation

- **Multilingual Support**
  - Full localization in 5 languages: English, Russian, Ukrainian, German, and Chinese
  - Language selection for both server and individual users
  - Easy addition of new languages through JSON localization files

- **YouTube Audio Playback**
  - Play audio from YouTube videos in voice channels
  - Simple playback control through commands
  - High sound quality thanks to optimized encoding

## 📋 Requirements

- Go 1.21 or higher
- SQLite3
- Discord Bot Token
- Gemini API Key (optional, for AI functions)
- Internet access for YouTube functions

## 🚀 Installation and Setup

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/discord-bot.git
cd discord-bot
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Create Configuration File

Create a `config.json` file in the root directory of the project:

```json
{
  "token": "YOUR_DISCORD_TOKEN",
  "prefix": "/",
  "gemini_api_key": "YOUR_GEMINI_API_KEY",
  "report_threshold": 3,
  "admin_role_id": "ADMIN_ROLE_ID",
  "mod_role_id": "MODERATOR_ROLE_ID",
  "default_language": "en",
  "bot_name": "Lapidar"
}
```

### 4. Compile and Run the Bot

```bash
go build
./discord-bot  # On Windows: discord-bot.exe
```

## 💬 Commands

### Basic Commands

| Command | Description | Access Rights |
|---------|-------------|---------------|
| `/help` | Show command help (displayed via webhook) | All users |
| `/report @user reason` | Report a user | All users |
| `/ban @user reason [duration]` | Ban a user | Administrators only |
| `/ai your query` | Ask a question to Gemini AI | All users |
| `/language [ru\|en\|uk\|de\|zh]` | Change bot language | All users |

### Voice Channel Commands

| Command | Description | Access Rights |
|---------|-------------|---------------|
| `/play YouTube-URL` | Play audio from YouTube | All users |
| `/stop` | Stop playback | All users |
| `/pause` | Pause playback | All users |
| `/resume` | Resume playback | All users |

## 🌐 Localization

The bot supports the following languages:
- English (en)
- Russian (ru)
- Ukrainian (uk)
- German (de)
- Chinese (Simplified) (zh)

Localization files are located in the `localization/translations/` directory.

### Switching Languages

You can switch the bot language in two ways:

1. **Server-wide setting**: Administrators can change the default language for the entire server using the command:
   ```
   /language [language_code]
   ```
   Where `language_code` is one of: `ru`, `en`, `uk`, `de`, `zh`

2. **User preference**: Individual users can set their preferred language using the same command, which will override the server default for them.

3. **Documentation**: This README is available in all supported languages. You can view it in your preferred language by opening the corresponding file in the `README/` directory:
   - [English](README/README.en.md)
   - [Russian](README/README.ru.md)
   - [Ukrainian](README/README.uk.md)
   - [German](README/README.de.md)
   - [Chinese](README/README.zh.md)

## 🎵 YouTube Audio Playback

To use the YouTube audio playback feature:
1. Join a voice channel
2. Enter the command `/play YouTube-URL`
3. The bot will join your voice channel and start playing audio
4. To stop playback, use the `/stop` command

## 🛠️ Report System

1. User sends a report using the `/report` command
2. The bot creates a report message in the moderation channel
3. Administrators can confirm or reject the report using reactions
4. When the report threshold is reached, the user is automatically banned

## 📁 Project Structure

- `main.go` - Main file for bot initialization
- `config/config.go` - Module for working with configuration
- `db/db.go` - Module for working with SQLite database
- `handlers/handlers.go` - Discord event handlers
- `handlers/gemini_handler.go` - Handler for Gemini AI integration
- `handlers/language_handler.go` - Handler for multilingual support
- `handlers/voice_handler.go` - Handler for voice functions
- `reports/reports.go` - Module for working with the report system
- `gemini/gemini.go` - Module for Gemini AI integration
- `localization/localization.go` - Module for localization
- `localization/translations/` - Translation files

## 📚 Used Libraries

- [discordgo](https://github.com/bwmarrin/discordgo) - Go package for Discord API
- [go-sqlite3](https://github.com/mattn/go-sqlite3) - SQLite3 driver for Go
- [google/generative-ai-go](https://github.com/google/generative-ai-go) - Google Generative AI client for Go
- [ytdl-core](https://github.com/fent/node-ytdl-core) - YouTube downloading module
- [dca](https://github.com/jonas747/dca) - Discord audio encoder

## 🔗 Connecting to Discord

1. Create a new application on the [Discord Developer Portal](https://discord.com/developers/applications)
2. Go to the "Bot" section and create a bot
3. Enable necessary Intents (Message Content, Server Members, Voice)
4. Copy the bot token and add it to your `config.json` file
5. Generate an invite link with required permissions in the OAuth2 section
6. Use the invite link to add the bot to your server

## 📜 License

[MIT License](LICENSE) © 2025 BANER

## 🔄 CI/CD и Автоматизация

Проект использует GitHub Actions для автоматизации проверки кода и создания релизов:

- **Проверка кода**: Автоматический линтинг и тестирование при каждом пуше и PR
- **Автоматические релизы**: При создании тега с версией (например, `v1.0.0`) автоматически создаётся релиз с бинарными файлами для Windows, Linux и macOS

Чтобы создать новый релиз:
1. Создайте и отправьте тег: `git tag v1.0.0 && git push origin v1.0.0`
2. GitHub Actions автоматически соберёт бинарные файлы и создаст релиз

## 🤝 Contributing

We welcome contributions to Lapidar! If you want to contribute:

1. Fork the repository
2. Create a branch for your feature (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to your fork (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📞 Support

If you have any issues or questions, please create an Issue in the project repository or contact us through Discord.

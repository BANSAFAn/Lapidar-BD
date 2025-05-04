# Lapidar - Discord Bot in Go

<div align="center">

<h1 style="border: none; font-size: 2.5em;">LAPIDAR</h1>

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go&logoColor=white)](https://golang.org/)
[![Discord](https://img.shields.io/badge/Discord-7289DA?logo=discord&logoColor=white)](https://discord.com/)
[![SQLite](https://img.shields.io/badge/SQLite-003B57?logo=sqlite&logoColor=white)](https://www.sqlite.org/)

</div>

<div align="center">

**Sprachversionen:** 
[English](README.en.md) | 
[Русский](README.ru.md) | 
[Українська](README.uk.md) | 
[Deutsch](README.de.md) | 
[中文](README.zh.md)

</div>

Lapidar ist ein leistungsstarker multifunktionaler Discord-Bot, der in Go geschrieben wurde. Er bietet ein breites Spektrum an Servermanagement-Funktionen, darunter ein Moderationssystem, KI-Integration, mehrsprachige Unterstützung und YouTube-Audiowiedergabe in Sprachkanälen.

## 🌟 Hauptfunktionen

- **Erweitertes Meldesystem**
  - Einreichen von Meldungen über Benutzer mit Administratorbestätigung
  - Automatische Sperrung beim Erreichen des Meldeschwellenwerts
  - Detaillierte Meldehistorie für jeden Benutzer

- **Moderationswerkzeuge**
  - Befehle zum Sperren von Benutzern mit Grund und Dauer
  - Flexibles Rollensystem für Administratoren und Moderatoren
  - Protokollierung aller Moderationsaktionen

- **KI-Integration**
  - Mehrere KI-Modelle (Gemini, Grok, ChatGPT, Qwen, Claude)
  - Stellen Sie Fragen an KI direkt in Discord
  - Erhalten Sie informative und genaue Antworten auf Ihre Anfragen
  - Anpassbare Parameter für die Antwortgenerierung

- **Mehrsprachige Unterstützung**
  - Vollständige Lokalisierung in 5 Sprachen: Englisch, Russisch, Ukrainisch, Deutsch und Chinesisch
  - Sprachauswahl sowohl für Server als auch für einzelne Benutzer
  - Einfaches Hinzufügen neuer Sprachen durch JSON-Lokalisierungsdateien

- **YouTube-Audiowiedergabe**
  - Abspielen von Audio aus YouTube-Videos in Sprachkanälen
  - Einfache Wiedergabesteuerung durch Befehle
  - Hohe Klangqualität dank optimierter Kodierung

## 📋 Anforderungen

- Go 1.21 oder höher
- SQLite3
- Discord Bot Token
- API-Schlüssel für KI-Modelle (optional)
- Internetzugang für YouTube-Funktionen

## 🚀 Installation und Einrichtung

### 1. Repository klonen

```bash
git clone https://github.com/yourusername/discord-bot.git
cd discord-bot
```

### 2. Abhängigkeiten installieren

```bash
go mod download
```

### 3. Konfigurationsdatei erstellen

Erstellen Sie eine `config.json` Datei im Stammverzeichnis des Projekts:

```json
{
  "token": "DEIN_DISCORD_TOKEN",
  "prefix": "!",
  "gemini_api_key": "DEIN_GEMINI_API_KEY",
  "grok_api_key": "DEIN_GROK_API_KEY",
  "chatgpt_api_key": "DEIN_CHATGPT_API_KEY",
  "qwen_api_key": "DEIN_QWEN_API_KEY",
  "claude_api_key": "DEIN_CLAUDE_API_KEY",
  "default_ai": "gemini",
  "report_threshold": 3,
  "admin_role_id": "ADMIN_ROLLEN_ID",
  "mod_role_id": "MODERATOR_ROLLEN_ID",
  "default_language": "de",
  "bot_name": "Lapidar"
}
```

### 4. Bot kompilieren und ausführen

```bash
go build
./discord-bot  # Unter Windows: discord-bot.exe
```

## 💬 Befehle

### Textbefehle

| Befehl | Beschreibung | Zugriffsrechte |
|--------|--------------|----------------|
| `!help` | Befehlshilfe anzeigen (wird über Webhook dargestellt) | Alle Benutzer |
| `!report @Benutzer Grund` | Einen Benutzer melden | Alle Benutzer |
| `!ban @Benutzer Grund [Dauer]` | Einen Benutzer sperren | Nur Administratoren |
| `!ai [Modell] deine Frage` | Eine Frage an KI stellen (mit Standard- oder angegebenem Modell) | Alle Benutzer |
| `!gemini deine Frage` | Eine Frage an Gemini AI stellen | Alle Benutzer |
| `!grok deine Frage` | Eine Frage an Grok AI stellen | Alle Benutzer |
| `!chatgpt deine Frage` | Eine Frage an ChatGPT stellen | Alle Benutzer |
| `!qwen deine Frage` | Eine Frage an Qwen AI stellen | Alle Benutzer |
| `!claude deine Frage` | Eine Frage an Claude AI stellen | Alle Benutzer |
| `!language [ru|en|uk|de|zh]` | Bot-Sprache ändern | Alle Benutzer |

### Sprachkanalbefehle

| Befehl | Beschreibung | Zugriffsrechte |
|--------|--------------|----------------|
| `!play YouTube-URL` | Audio von YouTube abspielen | Alle Benutzer |
| `!stop` | Wiedergabe stoppen | Alle Benutzer |
| `!pause` | Wiedergabe pausieren | Alle Benutzer |
| `!resume` | Wiedergabe fortsetzen | Alle Benutzer |

### Discord-Anwendungsbefehle (Slash-Befehle)

- `/ai` - Eine Frage an das Standard-KI-Modell stellen
- `/gemini` - Eine Frage an Gemini AI stellen
- `/grok` - Eine Frage an Grok AI stellen
- `/chatgpt` - Eine Frage an ChatGPT stellen
- `/qwen` - Eine Frage an Qwen AI stellen
- `/claude` - Eine Frage an Claude AI stellen

## 🌐 Lokalisierung

Der Bot unterstützt die folgenden Sprachen:
- Englisch (en)
- Russisch (ru)
- Ukrainisch (uk)
- Deutsch (de)
- Chinesisch (Vereinfacht) (zh)

Lokalisierungsdateien befinden sich im Verzeichnis `localization/translations/`.

### Sprachen wechseln

Sie können die Bot-Sprache auf zwei Arten umschalten:

1. **Serverweite Einstellung**: Administratoren können die Standardsprache für den gesamten Server mit dem Befehl ändern:
   ```
   !language [Sprachcode]
   ```
   Wobei `Sprachcode` einer der folgenden ist: `ru`, `en`, `uk`, `de`, `zh`

2. **Benutzereinstellung**: Einzelne Benutzer können ihre bevorzugte Sprache mit demselben Befehl festlegen, was die Servereinstellung für sie überschreibt.

## 🎵 YouTube-Audiowiedergabe

Um die YouTube-Audiowiedergabefunktion zu nutzen:
1. Einem Sprachkanal beitreten
2. Den Befehl `!play YouTube-URL` eingeben
3. Der Bot tritt Ihrem Sprachkanal bei und beginnt mit der Audiowiedergabe
4. Um die Wiedergabe zu stoppen, den Befehl `!stop` verwenden

## 🛠️ Meldesystem

1. Benutzer sendet eine Meldung mit dem Befehl `!report`
2. Der Bot erstellt eine Meldungsnachricht im Moderationskanal
3. Administratoren können die Meldung mit Reaktionen bestätigen oder ablehnen
4. Wenn der Meldeschwellenwert erreicht ist, wird der Benutzer automatisch gesperrt

## 📁 Projektstruktur

- `main.go` - Hauptdatei für Bot-Initialisierung
- `config/config.go` - Modul für die Arbeit mit der Konfiguration
- `db/db.go` - Modul für die Arbeit mit der SQLite-Datenbank
- `handlers/handlers.go` - Discord-Ereignishandler
- `reports/reports.go` - Modul für die Arbeit mit dem Meldesystem
- `ai/ai.go` - Modul für KI-Integration

## 📚 Verwendete Bibliotheken

- [discordgo](https://github.com/bwmarrin/discordgo) - Go-Paket für Discord API
- [go-sqlite3](https://github.com/mattn/go-sqlite3) - SQLite3-Treiber für Go
- [google/generative-ai-go](https://github.com/google/generative-ai-go) - Google Generative AI-Client für Go
- [ytdl-core](https://github.com/fent/node-ytdl-core) - YouTube-Download-Modul
- [dca](https://github.com/jonas747/dca) - Discord-Audio-Encoder

## 🔗 Verbindung mit Discord

1. Erstellen Sie eine neue Anwendung im [Discord Developer Portal](https://discord.com/developers/applications)
2. Gehen Sie zum Abschnitt "Bot" und erstellen Sie einen Bot
3. Aktivieren Sie die notwendigen Intents (Message Content, Server Members, Voice)
4. Kopieren Sie den Bot-Token und fügen Sie ihn in Ihre `config.json` Datei ein
5. Generieren Sie einen Einladungslink mit erforderlichen Berechtigungen im OAuth2-Abschnitt
6. Verwenden Sie den Einladungslink, um den Bot zu Ihrem Server hinzuzufügen

## 📄 Lizenz

MIT
# Lapidar - 基于 Go 的 Discord 机器人

<div align="center">

<h1 style="border: none; font-size: 2.5em;">LAPIDAR</h1>

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go&logoColor=white)](https://golang.org/)
[![Discord](https://img.shields.io/badge/Discord-7289DA?logo=discord&logoColor=white)](https://discord.com/)
[![SQLite](https://img.shields.io/badge/SQLite-003B57?logo=sqlite&logoColor=white)](https://www.sqlite.org/)

</div>

<div align="center">

**语言版本：** 
[English](README.en.md) | 
[Русский](README.ru.md) | 
[Українська](README.uk.md) | 
[Deutsch](README.de.md) | 
[中文](README.zh.md)

</div>

Lapidar 是一个用 Go 语言编写的功能强大的多功能 Discord 机器人。它提供了广泛的服务器管理功能，包括审核系统、AI 集成、多语言支持以及在语音频道中播放 YouTube 音频。

## 🌟 主要功能

- **高级举报系统**
  - 提交用户举报并由管理员确认
  - 达到举报阈值时自动封禁
  - 每个用户的详细举报历史

- **审核工具**
  - 带有原因和时长的用户封禁命令
  - 管理员和审核员的灵活角色系统
  - 记录所有审核操作

- **AI 集成**
  - 多种 AI 模型（Gemini、Grok、ChatGPT、Qwen、Claude）
  - 直接在 Discord 中向 AI 提问
  - 获取信息丰富且准确的回答
  - 可自定义的回答生成参数

- **多语言支持**
  - 5 种语言的完整本地化：英语、俄语、乌克兰语、德语和中文
  - 服务器和个人用户的语言选择
  - 通过 JSON 本地化文件轻松添加新语言

- **YouTube 音频播放**
  - 在语音频道中播放 YouTube 视频的音频
  - 通过命令简单控制播放
  - 优化编码带来的高音质

## 📋 要求

- Go 1.21 或更高版本
- SQLite3
- Discord Bot Token
- AI 模型的 API 密钥（可选）
- YouTube 功能需要互联网访问

## 🚀 安装和设置

### 1. 克隆仓库

```bash
git clone https://github.com/yourusername/discord-bot.git
cd discord-bot
```

### 2. 安装依赖

```bash
go mod download
```

### 3. 创建配置文件

在项目根目录创建 `config.json` 文件：

```json
{
  "token": "您的_DISCORD_TOKEN",
  "prefix": "!",
  "gemini_api_key": "您的_GEMINI_API_KEY",
  "grok_api_key": "您的_GROK_API_KEY",
  "chatgpt_api_key": "您的_CHATGPT_API_KEY",
  "qwen_api_key": "您的_QWEN_API_KEY",
  "claude_api_key": "您的_CLAUDE_API_KEY",
  "default_ai": "gemini",
  "report_threshold": 3,
  "admin_role_id": "管理员角色ID",
  "mod_role_id": "审核员角色ID",
  "default_language": "zh",
  "bot_name": "Lapidar"
}
```

### 4. 编译并运行机器人

```bash
go build
./discord-bot  # 在 Windows 上：discord-bot.exe
```

## 💬 命令

### 文本命令

| 命令 | 描述 | 访问权限 |
|-----|------|--------|
| `!help` | 显示命令帮助（通过 webhook 显示） | 所有用户 |
| `!report @用户 原因` | 举报用户 | 所有用户 |
| `!ban @用户 原因 [时长]` | 封禁用户 | 仅管理员 |
| `!ai [模型] 您的问题` | 向 AI 提问（使用默认或指定的模型） | 所有用户 |
| `!gemini 您的问题` | 向 Gemini AI 提问 | 所有用户 |
| `!grok 您的问题` | 向 Grok AI 提问 | 所有用户 |
| `!chatgpt 您的问题` | 向 ChatGPT 提问 | 所有用户 |
| `!qwen 您的问题` | 向 Qwen AI 提问 | 所有用户 |
| `!claude 您的问题` | 向 Claude AI 提问 | 所有用户 |
| `!language [ru|en|uk|de|zh]` | 更改机器人语言 | 所有用户 |

### 语音频道命令

| 命令 | 描述 | 访问权限 |
|-----|------|--------|
| `!play YouTube-URL` | 播放 YouTube 音频 | 所有用户 |
| `!stop` | 停止播放 | 所有用户 |
| `!pause` | 暂停播放 | 所有用户 |
| `!resume` | 恢复播放 | 所有用户 |

### Discord 应用命令（斜杠命令）

- `/ai` - 向默认 AI 模型提问
- `/gemini` - 向 Gemini AI 提问
- `/grok` - 向 Grok AI 提问
- `/chatgpt` - 向 ChatGPT 提问
- `/qwen` - 向 Qwen AI 提问
- `/claude` - 向 Claude AI 提问

## 🌐 本地化

机器人支持以下语言：
- 英语 (en)
- 俄语 (ru)
- 乌克兰语 (uk)
- 德语 (de)
- 中文（简体）(zh)

本地化文件位于 `localization/translations/` 目录中。

### 切换语言

您可以通过两种方式切换机器人语言：

1. **服务器范围设置**：管理员可以使用以下命令更改整个服务器的默认语言：
   ```
   !language [语言代码]
   ```
   其中 `语言代码` 是以下之一：`ru`、`en`、`uk`、`de`、`zh`

2. **用户偏好**：个人用户可以使用相同的命令设置他们的首选语言，这将覆盖服务器对他们的设置。

## 🎵 YouTube 音频播放

要使用 YouTube 音频播放功能：
1. 加入语音频道
2. 输入命令 `!play YouTube-URL`
3. 机器人将加入您的语音频道并开始播放音频
4. 要停止播放，使用 `!stop` 命令

## 🛠️ 举报系统

1. 用户使用 `!report` 命令发送举报
2. 机器人在审核频道创建举报消息
3. 管理员可以通过反应确认或拒绝举报
4. 当达到举报阈值时，用户将被自动封禁

## 📁 项目结构

- `main.go` - 机器人初始化的主文件
- `config/config.go` - 配置工作模块
- `db/db.go` - SQLite 数据库工作模块
- `handlers/handlers.go` - Discord 事件处理程序
- `reports/reports.go` - 举报系统工作模块
- `ai/ai.go` - AI 集成模块

## 📚 使用的库

- [discordgo](https://github.com/bwmarrin/discordgo) - Go 语言的 Discord API 包
- [go-sqlite3](https://github.com/mattn/go-sqlite3) - Go 语言的 SQLite3 驱动
- [google/generative-ai-go](https://github.com/google/generative-ai-go) - Go 语言的 Google Generative AI 客户端
- [ytdl-core](https://github.com/fent/node-ytdl-core) - YouTube 下载模块
- [dca](https://github.com/jonas747/dca) - Discord 音频编码器

## 🔗 连接到 Discord

1. 在 [Discord Developer Portal](https://discord.com/developers/applications) 创建新应用
2. 转到 "Bot" 部分并创建机器人
3. 启用必要的 Intents（Message Content, Server Members, Voice）
4. 复制机器人令牌并添加到您的 `config.json` 文件中
5. 在 OAuth2 部分生成具有所需权限的邀请链接
6. 使用邀请链接将机器人添加到您的服务器

## 📄 许可证

MIT
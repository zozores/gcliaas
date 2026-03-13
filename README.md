# GCLIaaS (Gemini CLI as a Server)

This repository contains two interconnected Go projects that expose the existing `gemini` CLI tool as an HTTP service and consume it via a Telegram bot.

## Requirements

- `gemini` CLI tool
- `systemd` user service
- `go` 1.26.1

## Components

### 1. `gemini_listen`
A lightweight Go HTTP server that acts as a bridge between HTTP requests and the `gemini` CLI tool (`gemini --yolo --prompt <message> -m gemini-2.5-flash`).

- **Endpoint**: `POST /event`
- **Payload Format**: `{"message": "your prompt here"}`
- **Default Port**: `8765`

#### Setup
```bash
cd gemini_listen
make install # Builds the binary, copies it to ~/.local/bin, and sets up a systemd user service.
```

### 2. `telegram_bot`
A Telegram bot written in Go that forwards user messages (text or voice) to the `gemini_listen` service.

- Supports both text messages and voice transcription.
- Voice transcription requires a native Gemini API key (can be set via `gemini-telegram-bot.service` or during the chat by just passing it).
- **Environment Variables** (see `gemini-telegram-bot.service`):
  - `TELEGRAM_BOT_TOKEN`: Your Telegram bot token. (Required)
  - `GEMINI_ENDPOINT`: URL of the `gemini_listen` service. (Default: `http://127.0.0.1:8765/event`)
  - `GEMINI_API_KEY`: Required for voice transcription of audio messages.
  - `TARGET_CHAT_ID`: (Optional) Restrict the bot to only respond to a specific chat ID.

#### Setup
```bash
cd telegram_bot
make install # Builds the binary, copies it to ~/.local/bin, and sets up a systemd user service.
```

## Running

Both services include `Makefile`s to build, install, and manage them as `systemd` user services. 

To manage them manually, check `systemctl --user status gemini-listen.service` and `systemctl --user status gemini-telegram-bot.service`. You can use `make status`, `make restart`, `make uninstall` to manage their lifecycle from within their respective directories.

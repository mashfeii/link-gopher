# Telegram Link Tracker Bot

A service for tracking web resource updates via Telegram, with bot and scrapper components communicating over HTTP.

## Requirements

- **[Go 1.23+](github.com/golang/go)**
- **[Docker 24.0+](docker.com)** (optional for containerized setup)
- **Telegram Bot Token** (from [@BotFather](https://t.me/BotFather))

## Installation

```bash
git clone git@github.com:central-university-dev/go-mashfeii.git
cd go-mashfeii
```

## Configuration

Default configuration lies in `config/default-config.yaml`. You can override these values by specifying custom
configuration file path within `--config` flag on execution binary for both bot and scrapper.

Project contains several secrets that should be provided via environment, the only mandatory one is `BOT_TOKEN`, without
`GITHUB_TOKEN` or `STACKOVERFLOW_TOKEN` you will be limited in number of requests to respective APIs.

### Local Development

```bash
# Build binaries into `bin` directory
make build

# Start services
BOT_TOKEN={your_token} ./bin/bot&
GITHUB_TOKEN={your_token} STACKOVERFLOW_TOKEN={your_token} ./bin/scrapper&
```

### Docker

```bash
# Export secrets
export BOT_TOKEN={your_token}
export GITHUB_TOKEN={your_token}
export STACKOVERFLOW_TOKEN={your_token}

# Build and run containers
docker-compose -f docker/docker-compose.yaml up --build
```

## Usage

### Telegram Commands

| Command    | Description                  |
| ---------- | ---------------------------- |
| `/start`   | Register new user in service |
| `/track`   | Start monitoring a link      |
| `/untrack` | Stop monitoring a link       |
| `/list`    | List all tracked links       |
| `/help`    | Display commands list        |
| `/cancel`  | Cancel current operation     |

### API Examples

#### Scrapper API

```bash
# Register chat
curl -X POST http://localhost:8080/tg-chat/{chat_id}

# List links (requires chat id inside header)
curl -H "Tg-Chat-Id: {chat_id}" http://localhost:8080/links
```

#### Bot API

```bash
# Send update notification
curl -X POST http://localhost:8081/updates \
  -H "Content-Type: application/json" \
  -d '{
    "url": "stackoverflow.com/questions/292357",
    "description": "New answer added",
    "tgChatId": 7517423324
  }'
```

### Project Structure

```
├── cmd/                # Entry points (bot/scrapper)
├── config/             # Configuration loader
├── docker/             # Docker configurations
├── api/                # OpenAPI specifications
├── internal/           # Core application logic
│   ├── clients/        # Third-party integrations
│   ├── domain/         # Business logic models
│   └── infrastructure/ # Implementations
└── Makefile            # Build automation
```

### Development Tools

#### Prerequisites

- [golangci-lint](github.com/golangci/golangci-lint): linters aggregator
- [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen/): openAPI generator

```bash
make test        # Run tests with coverage report
make lint        # Check code quality (Go + Proto)
make generate    # Generate API/client code
make clean       # Remove build artifacts
```

### OpenAPI Specifications

- **[Scrapper API](./api/openapi/v1/scrapper-api.yaml)**
- **[Bot API](./api/openapi/v1/bot-api.yaml)**

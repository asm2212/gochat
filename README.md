# gochat

Professional, scalable, and clean Go backend for a Telegram-style chat system:

- User signup & login (JWT, bcrypt)
- Direct one-to-one messages
- Group chat (create, send, fetch)
- Global broadcast messaging

## Project Structure
cmd/server/main.go # Entry point internal/api.go # HTTP handlers internal/auth.go # JWT middleware internal/user.go # User logic internal/chat.go # Messaging logic internal/model.go # Data models internal/redis.go # Redis connection go.mod 


## Quick Start

1. Start Redis:
go mod tidy go run ./cmd/server/main.go

3. API endpoints:
- `POST   /signup`  `{ "username": "...", "password": "..." }`
- `POST   /login`   `{ "username": "...", "password": "..." }` → `{ "token": "..." }`
- `POST   /dm/send` `{ "to": "...", "content": "..." }` (auth)
- `GET    /dm/history?user=...` (auth)
- `POST   /group/create` `{ "group": "..." }` (auth)
- `POST   /group/send` `{ "group": "...", "content": "..." }` (auth)
- `GET    /group/history?group=...` (auth)
- `POST   /broadcast/send` `{ "content": "..." }` (auth)
- `GET    /broadcast/history` (auth)

Use JWT `Authorization: Bearer ...` for all but signup/login.

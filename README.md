# gochat

A professional, scalable, and clean Go backend for a Telegram-style chat system.

**Features:**
- User signup & login (JWT, bcrypt)
- Direct one-to-one messages
- Group chat (create, send, fetch)
- Global broadcast messaging

---

## 📁 Project Structure

```
gochat/
  main.go             # Entry point
  internal/
    api.go            # HTTP handlers
    auth.go           # JWT middleware
    user.go           # User logic
    chat.go           # Messaging logic
    model.go          # Data models
    redis.go          # Redis connection
  go.mod
  README.md
```

---

## 🚀 Quick Start

1. **Clone the repository**
    ```sh
    git clone https://github.com/asm2212/gochat.git
    cd gochat
    ```

2. **Start a Redis server**  
   You must have Redis running on your local machine (`localhost:6379`).
   - You can [download and run Redis](https://redis.io/download/) for your operating system.

3. **Install Go dependencies**
    ```sh
    go mod tidy
    ```

4. **Run the server**
    ```sh
    go run .
    ```
    The server will start on port `:8080`.

---

## 📖 API Endpoints

> All endpoints except `/signup` and `/login` require JWT authentication.  
> Provide the token in the `Authorization: Bearer <token>` HTTP header.

### Authentication

- **Sign Up**
  - `POST /signup`
  - Body:
    ```json
    {
      "username": "your_username",
      "password": "your_password"
    }
    ```
- **Login**
  - `POST /login`
  - Body:
    ```json
    {
      "username": "your_username",
      "password": "your_password"
    }
    ```
  - Response:
    ```json
    {
      "token": "<JWT_TOKEN>"
    }
    ```

### Direct Messages

- **Send a direct message**
  - `POST /dm/send`
  - Header: `Authorization: Bearer <token>`
  - Body:
    ```json
    {
      "to": "recipient_username",
      "content": "Hello!"
    }
    ```
- **Fetch direct message history**
  - `GET /dm/history?user=recipient_username`
  - Header: `Authorization: Bearer <token>`

### Group Chat

- **Create a group**
  - `POST /group/create`
  - Header: `Authorization: Bearer <token>`
  - Body:
    ```json
    {
      "group": "group_name"
    }
    ```

- **Send a group message**
  - `POST /group/send`
  - Header: `Authorization: Bearer <token>`
  - Body:
    ```json
    {
      "group": "group_name",
      "content": "Group message!"
    }
    ```

- **Fetch group chat history**
  - `GET /group/history?group=group_name`
  - Header: `Authorization: Bearer <token>`

### Broadcast

- **Send a broadcast message**
  - `POST /broadcast/send`
  - Header: `Authorization: Bearer <token>`
  - Body:
    ```json
    {
      "content": "Message to everyone!"
    }
    ```
- **Fetch broadcast messages**
  - `GET /broadcast/history`
  - Header: `Authorization: Bearer <token>`

---


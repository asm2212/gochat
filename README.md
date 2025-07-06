# gochat

A Go backend for a Telegram-style chat system.

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

_All endpoints except `/signup` and `/login` require JWT authentication.  
Provide the token in the `Authorization: Bearer <token>` HTTP header._

---

### 🔒 Authentication

#### Sign Up

- **Endpoint:** `POST /signup`
- **Request Body:**
    ```json
    {
      "username": "your_username",
      "password": "your_password"
    }
    ```
- **Response:**
    - `201 Created` on success:  
      ```json
      { "message": "user registered" }
      ```
    - `409 Conflict` if user exists:  
      ```json
      { "error": "username already exists" }
      ```

#### Login

- **Endpoint:** `POST /login`
- **Request Body:**
    ```json
    {
      "username": "your_username",
      "password": "your_password"
    }
    ```
- **Response:**
    - `200 OK` on success:  
      ```json
      { "token": "<JWT_TOKEN>" }
      ```
    - `401 Unauthorized` on failure:  
      ```json
      { "error": "invalid credentials" }
      ```

---

### 💬 Direct Messages

#### Send a Direct Message

- **Endpoint:** `POST /dm/send`
- **Headers:**  
  `Authorization: Bearer <token>`
- **Request Body:**
    ```json
    {
      "to": "recipient_username",
      "content": "Hello!"
    }
    ```
- **Response:**
    - `200 OK`  
      ```json
      { "message": "sent" }
      ```
    - `500 Internal Server Error`  
      ```json
      { "error": "..." }
      ```

#### Fetch Direct Message History

- **Endpoint:** `GET /dm/history?user=recipient_username`
- **Headers:**  
  `Authorization: Bearer <token>`
- **Response:**
    - `200 OK`  
      ```json
      [
        {
          "id": "b7c9...",
          "from": "alice",
          "to": "bob",
          "content": "Hello Bob!",
          "type": "direct",
          "timestamp": "2025-07-06T05:00:00Z"
        },
        ...
      ]
      ```
    - `500 Internal Server Error`  
      ```json
      { "error": "..." }
      ```

---

### 👥 Group Chat

#### Create a Group

- **Endpoint:** `POST /group/create`
- **Headers:**  
  `Authorization: Bearer <token>`
- **Request Body:**
    ```json
    {
      "group": "group_name"
    }
    ```
- **Response:**
    - `201 Created`  
      ```json
      { "message": "group created" }
      ```
    - `500 Internal Server Error`  
      ```json
      { "error": "..." }
      ```

#### Send a Group Message

- **Endpoint:** `POST /group/send`
- **Headers:**  
  `Authorization: Bearer <token>`
- **Request Body:**
    ```json
    {
      "group": "group_name",
      "content": "Group message!"
    }
    ```
- **Response:**
    - `200 OK`  
      ```json
      { "message": "sent" }
      ```
    - `500 Internal Server Error`  
      ```json
      { "error": "..." }
      ```

#### Fetch Group Chat History

- **Endpoint:** `GET /group/history?group=group_name`
- **Headers:**  
  `Authorization: Bearer <token>`
- **Response:**
    - `200 OK`  
      ```json
      [
        {
          "id": "a1b2...",
          "from": "alice",
          "group": "mygroup",
          "content": "Hi everyone",
          "type": "group",
          "timestamp": "2025-07-06T05:00:00Z"
        },
        ...
      ]
      ```
    - `500 Internal Server Error`  
      ```json
      { "error": "..." }
      ```

---

### 📢 Broadcast

#### Send a Broadcast Message

- **Endpoint:** `POST /broadcast/send`
- **Headers:**  
  `Authorization: Bearer <token>`
- **Request Body:**
    ```json
    {
      "content": "Message to everyone!"
    }
    ```
- **Response:**
    - `200 OK`  
      ```json
      { "message": "broadcasted" }
      ```
    - `500 Internal Server Error`  
      ```json
      { "error": "..." }
      ```

#### Fetch Broadcast Messages

- **Endpoint:** `GET /broadcast/history`
- **Headers:**  
  `Authorization: Bearer <token>`
- **Response:**
    - `200 OK`  
      ```json
      [
        {
          "id": "xxxx",
          "from": "alice",
          "content": "Hello all!",
          "type": "broadcast",
          "timestamp": "2025-07-06T05:00:00Z"
        },
        ...
      ]
      ```
    - `500 Internal Server Error`  
      ```json
      { "error": "..." }
      ```

---


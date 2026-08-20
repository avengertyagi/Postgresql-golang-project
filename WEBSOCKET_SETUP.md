# WebSocket Setup Guide

This guide will help you set up and use WebSocket functionality in your PostgreSQL-Golang project.

## 📋 Overview

The WebSocket module is built with:
- **Gorilla WebSocket**: Reliable WebSocket library
- **Hub Architecture**: Manages multiple connections and broadcasting
- **Client Model**: Handles individual connections with read/write pumps
- **Gin Integration**: Seamless integration with your existing Gin router

## 🚀 Installation

### 1. Add gorilla/websocket Dependency

```bash
go get github.com/gorilla/websocket
go mod tidy
```

### 2. Verify Files Are Created

The following files should now be in place:

```
internal/modules/websocket/
├── controllers/
│   └── websocket-controller.go
├── models/
│   └── message.go
├── routes/
│   └── websocket-route.go
└── services/
    └── chat.go
```

## 🔧 Configuration

### Environment Variables

Add these to your `.env` file if needed:

```env
# WebSocket Configuration
WS_READ_BUFFER_SIZE=1024
WS_WRITE_BUFFER_SIZE=1024
WS_PONG_WAIT=60
WS_PING_PERIOD=54
```

### CORS Setup

Update your CORS configuration in `cmd/api/main.go` to allow WebSocket:

```go
cors.New(cors.Config{
    AllowOrigins:     strings.Split(appConfig.AllowedOrigin, ","),
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true, // Important for WebSocket
    MaxAge:           12 * time.Hour,
})
```

## 🏃 Running the Application

1. **Start the server:**
   ```bash
   make dev
   # or
   go run ./cmd/api
   ```

2. **The WebSocket server will be available at:**
   ```
   ws://localhost:8080/api/v1/ws
   ```

3. **Check stats at:**
   ```
   http://localhost:8080/api/v1/ws/stats
   ```

## 📡 API Endpoints

### WebSocket Endpoint
- **URL**: `GET /ws`
- **Protocol**: WebSocket
- **Authentication**: Required (Bearer token in query or header)
- **Description**: Main WebSocket connection endpoint

### Stats Endpoint
- **URL**: `GET /ws/stats`
- **Method**: HTTP GET
- **Response**: JSON with connected client count

## 💬 Client Implementation

### JavaScript Client Example

```javascript
// Connect to WebSocket
const ws = new WebSocket('ws://localhost:8080/api/v1/ws?token=YOUR_JWT_TOKEN');

// Connection opened
ws.onopen = (event) => {
  console.log('WebSocket connected');
  ws.send(JSON.stringify({
    type: 'text',
    content: 'Hello from client',
    timestamp: new Date().toISOString()
  }));
};

// Receive message
ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  console.log('Received:', message);
};

// Connection error
ws.onerror = (error) => {
  console.error('WebSocket error:', error);
};

// Connection closed
ws.onclose = (event) => {
  console.log('WebSocket closed');
};
```

### Go Client Example

```go
package main

import (
    "fmt"
    "github.com/gorilla/websocket"
    "log"
)

func main() {
    url := "ws://localhost:8080/api/v1/ws"
    ws, _, err := websocket.DefaultDialer.Dial(url, nil)
    if err != nil {
        log.Fatal(err)
    }
    defer ws.Close()

    // Send message
    err = ws.WriteMessage(websocket.TextMessage, []byte("Hello"))
    if err != nil {
        log.Fatal(err)
    }

    // Receive message
    _, message, err := ws.ReadMessage()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Received: %s\n", message)
}
```

## 🏗️ Architecture

### Hub
- Manages all connected clients
- Handles registration and unregistration
- Broadcasts messages to all clients
- Runs in its own goroutine

### Client
- Represents a single WebSocket connection
- Has separate read and write pumps
- Maintains a send channel for outgoing messages
- Handles ping/pong for connection keep-alive

### Message Flow
```
Client A sends → Hub receives → Hub broadcasts → All clients receive
```

## 📊 Message Format

### Request Message
```json
{
  "id": "msg-123",
  "type": "text",
  "content": "Your message here",
  "timestamp": "2026-08-20T10:30:00Z",
  "data": {}
}
```

### Response Message
```json
{
  "id": "msg-123",
  "type": "text",
  "content": "Your message here",
  "user_id": "user-1",
  "user_name": "John Doe",
  "timestamp": "2026-08-20T10:30:00Z"
}
```

## 🔐 Authentication

### Token in Query Parameter
```
ws://localhost:8080/api/v1/ws?token=YOUR_JWT_TOKEN
```

### Token in Header (Custom Implementation)
Update the controller to extract from headers:
```go
func (wsc *WebSocketController) HandleWebSocket(c *gin.Context) {
    authHeader := c.GetHeader("Authorization")
    // Parse and validate token
}
```

## 🧪 Testing

### Manual Testing with wscat

```bash
# Install wscat
npm install -g wscat

# Connect to WebSocket
wscat -c ws://localhost:8080/api/v1/ws

# Send messages
> {"type": "text", "content": "Hello"}
< {"type": "text", "content": "Hello", "user_id": "1", ...}
```

### Using cURL

```bash
# Check stats endpoint
curl http://localhost:8080/api/v1/ws/stats
```

## 🐛 Troubleshooting

### Connection Refused
- Ensure server is running
- Check port 8080 is available
- Verify CORS configuration

### Authentication Failed
- Verify JWT token is valid
- Check token format in query parameter
- Ensure auth middleware is properly configured

### Messages Not Broadcasting
- Check hub is running: `go websocketHub.Run()` in container
- Verify client send channels are not full
- Check message size limits (max 512 bytes by default)

### High Memory Usage
- Monitor connected clients with `/ws/stats`
- Increase read/write timeouts if clients are slow
- Consider implementing message rate limiting

## 🔄 Scaling Considerations

### Multiple Instances
For multiple server instances, consider using:
- Redis for pub/sub across instances
- Message queue (RabbitMQ, Kafka) for distributed messaging
- Load balancer with sticky sessions

### Example with Redis (Future Enhancement)
```go
// Subscribe to Redis channel
pubsub := redisClient.Subscribe(ctx, "chat-channel")

// Broadcast through Redis
redisClient.Publish(ctx, "chat-channel", message)
```

## 📝 Best Practices

1. **Always validate messages** before broadcasting
2. **Set appropriate timeouts** for read/write operations
3. **Implement rate limiting** to prevent abuse
4. **Clean up disconnected clients** properly
5. **Monitor connection count** for performance
6. **Use structured logging** for debugging
7. **Handle errors gracefully** without crashing
8. **Test with multiple concurrent connections**

## 🚀 Advanced Features (To Implement)

1. **Rooms/Channels**: Group messages by topics
2. **Private Messages**: Direct messaging between users
3. **Message History**: Store and retrieve past messages
4. **Typing Indicators**: Show when users are typing
5. **User Presence**: Track online/offline status
6. **Message Persistence**: Store in database
7. **Notification System**: Push notifications
8. **Message Encryption**: End-to-end encryption

## 📚 Additional Resources

- [Gorilla WebSocket Documentation](https://pkg.go.dev/github.com/gorilla/websocket)
- [WebSocket Protocol (RFC 6455)](https://tools.ietf.org/html/rfc6455)
- [Gin Framework Documentation](https://gin-gonic.com/)

## ⚠️ Important Notes

- The current implementation allows all origins in `CheckOrigin`. Update for production:
  ```go
  CheckOrigin: func(r *http.Request) bool {
      origin := r.Header.Get("Origin")
      return isAllowedOrigin(origin)
  }
  ```

- Max message size is 512 bytes. Adjust in `chat.go` if needed:
  ```go
  const maxMessageLen = 512 // Change as needed
  ```

- Connection timeout is 60 seconds. Adjust based on your needs:
  ```go
  const pongWait = 60 * time.Second
  ```

---

**Last Updated**: 2026-08-20

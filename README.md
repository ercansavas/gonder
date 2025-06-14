# Gonder 🚀

Gonder is a modern messaging service written in Go. It provides a flexible API for sending messages via Email, SMS, and other communication channels.

## ✨ Features

- 🌐 RESTful API
- 📊 **Comprehensive Audit Logging** - All system events are written to console in JSON format
- 📧 Email sending support (planned)
- 📱 SMS sending support (planned)
- 🔧 Easy configuration
- 🐳 Docker support (planned)
- ⚡ High performance - Go 1.24.4
- 🛡️ Error handling and validation

## 📊 Audit Logging

The system comes with **comprehensive audit logging** feature:

### Logged Events:
- ✅ **API Calls** - All HTTP requests (method, path, status, duration, IP, user-agent)
- ✅ **Message Sent** - Message sending details (recipient, type, success/failure)
- ✅ **Errors** - System errors and validation errors
- ✅ **Startup/Shutdown** - Application lifecycle
- ✅ **Health Checks** - System health checks

### Log Format:
```json
{
  "timestamp": "2025-06-15T01:57:27.982286285+03:00",
  "event_type": "api_call",
  "method": "POST",
  "path": "/api/send",
  "status_code": 200,
  "duration": "1.234ms", 
  "message": "POST /api/send - 200",
  "details": {...},
  "remote_addr": "127.0.0.1:12345",
  "user_agent": "curl/7.68.0"
}
```

## 🚀 Installation

### Requirements
- Go 1.24.4 or higher
- Git

### Running

```bash
# Clone the project
git clone <repository-url>
cd gonder

# Install dependencies
go mod tidy

# Run the application
go run cmd/gonder/main.go
```

### Environment Variables
```bash
PORT=8080        # Server port (default: 8080)
HOST=localhost   # Server host (default: localhost) 
LOG_LEVEL=info   # Log level (default: info)
```

## 📋 API Endpoints

### Home Page
```
GET /
```
HTML homepage

### Send Message
```
POST /api/send
Content-Type: application/json

{
  "message": "Hello World!",
  "recipient": "user@example.com",
  "type": "email"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Message sent successfully",
  "id": "msg_1234567890",
  "timestamp": "2025-06-15T01:57:48+03:00"
}
```

### Health Check
```
GET /api/health
```

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2025-06-15T01:57:48+03:00",
  "version": "1.0.0",
  "uptime": "N/A"
}
```

## 🏗️ Project Structure

```
gonder/
├── cmd/gonder/main.go          # Main application
├── pkg/
│   ├── audit/                  # Audit logging system
│   │   ├── audit.go           # Audit logger and event types
│   │   └── middleware.go      # HTTP middleware
│   ├── handler/handler.go      # HTTP handlers
│   └── model/                  # Data models
├── internal/config/config.go   # Configuration
├── docs/                       # Documentation
├── go.mod                      # Go module
└── README.md                   # This file
```

## 🧪 Testing

```bash
# Health check
curl http://localhost:8080/api/health

# Send message
curl -X POST http://localhost:8080/api/send \
  -H "Content-Type: application/json" \
  -d '{"message":"Test message","recipient":"test@example.com"}'

# Error test (validation)
curl -X POST http://localhost:8080/api/send \
  -H "Content-Type: application/json" \
  -d '{"message":"","recipient":"test@example.com"}'
```

## 📈 Sample Audit Logs

### Application Startup
```json
[AUDIT] {"timestamp":"2025-06-15T01:57:27+03:00","event_type":"startup","message":"Gonder application started - Port: 8080","details":{"host":"localhost","log_level":"info","version":"1.0.0"}}
```

### API Call
```json
[AUDIT] {"timestamp":"2025-06-15T01:57:48+03:00","event_type":"api_call","method":"POST","path":"/api/send","status_code":200,"duration":"2.1ms","message":"POST /api/send - 200","details":{"content_type":"application/json","content_length":75},"remote_addr":"127.0.0.1:45678","user_agent":"curl/7.68.0"}
```

### Message Sending
```json
[AUDIT] {"timestamp":"2025-06-15T01:57:48+03:00","event_type":"message_sent","message":"Message sent: email -> demo@example.com (ID: msg_1749941868)","details":{"recipient":"demo@example.com","message_type":"email","message_id":"msg_1749941868","success":true,"extra":{"message_length":25,"message_preview":"Audit log demo message"}}}
```

### Error Case
```json
[AUDIT] {"timestamp":"2025-06-15T01:57:48+03:00","event_type":"error","message":"Error in Validation error in Send endpoint: message field is empty","error":"message field is empty","details":{"request":{"message":"","recipient":"demo@example.com","type":"email"}}}
```

## 🔧 Development

```bash
# Run tests
go test ./...

# Build
go build -o gonder cmd/gonder/main.go

# Format
go fmt ./...

# Vet
go vet ./...
```

## 📝 TODO

- [ ] Real email sending integration
- [ ] SMS sending support
- [ ] Database integration
- [ ] Authentication & authorization
- [ ] Rate limiting
- [ ] Metrics & monitoring
- [ ] Docker containerization
- [ ] CI/CD pipeline

## 📄 License

This project is licensed under the MIT License.
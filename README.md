# Gonder 🚀

**Gonder** is a modern system log collection service written in Go. It provides real-time log monitoring, parsing, and structured output with comprehensive audit logging.

## What is Gonder?

Gonder collects system logs, parses them in real-time, and outputs structured JSON format. It includes:

- 📊 **Real-time log collection** from system log files
- 🔍 **Automatic log parsing** with structured output
- 📋 **RESTful API** for management and monitoring
- 🛡️ **Comprehensive audit logging** for all system events
- ⚡ **High performance** with Go 1.24.4
- 🐳 **Docker support** for easy deployment

## 🚀 Quick Start

### With Docker (Recommended)

```bash
# Clone and run
git clone <repository-url>
cd gonder
docker-compose up -d

# Test it's working
curl http://localhost:8080/api/health
```

### Without Docker

```bash
# Clone the project
git clone <repository-url>
cd gonder

# Run
go mod tidy
go run cmd/gonder/main.go
```

**Access the service:** http://localhost:8080

## 📋 Main Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/` | GET | Homepage |
| `/api/health` | GET | Health check |
| `/api/logs/status` | GET | Log collector status |
| `/api/logs/sources` | GET | List log sources |
| `/api/logs/start` | POST | Start collector |
| `/api/logs/stop` | POST | Stop collector |

For complete API documentation, see [docs/api/API.md](docs/api/API.md).

## 📚 Documentation

Complete documentation is available in the **[docs/](docs/)** directory:

**🐳 Docker & Deployment:**
- [Quick Start Guide](docs/docker/QUICK-START.md)
- [Docker Guide](docs/docker/DOCKER.md)
- [Docker Commands](docs/docker/DOCKER-COMMANDS.md)

**📋 API & Architecture:**
- [API Documentation](docs/api/API.md)
- [Technical Design](docs/architecture/Technical-Design.md)
- [Product Requirements](docs/architecture/PRD.md)

👉 **See [docs/README.md](docs/README.md) for the complete documentation index.**

## 📄 License

This project is licensed under the MIT License.
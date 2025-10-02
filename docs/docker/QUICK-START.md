# 🚀 Gonder Quick Start Guide

## 📦 Project Overview

**Gonder v2.0** is a modern **system log collection service** written in Go.

### Key Features:
- 📊 Real-time system log collection
- 🔍 Automatic log parsing with structured output
- 📋 RESTful API management
- 🛡️ Comprehensive audit logging
- 🐳 Docker support

### How It Works:
The application reads system log files (e.g., `/var/log/system.log`), parses them, and outputs structured JSON to the console. All operations are recorded by the audit logging system.

---

## 🐳 Running with Docker (Recommended)

### 1️⃣ Simple Method - Docker Compose

```bash
# Navigate to project directory
cd /Users/ercan/Desktop/Projects/GitHub/gonder

# Build and run
docker-compose up -d

# Follow logs
docker-compose logs -f gonder
```

**Expected output:**
```
🚀 Gonder - System Log Collection Service starting...
🔧 Starting system log collector...
✅ System log collector started successfully
🌐 Server running on port 8080
```

### 2️⃣ Manual Docker Commands

```bash
# 1. Build the image
docker build -t gonder:latest .

# 2. Start the container
docker run -d \
  --name gonder \
  -p 8080:8080 \
  -e PORT=8080 \
  -e HOST=0.0.0.0 \
  -e LOG_LEVEL=info \
  gonder:latest

# 3. View logs
docker logs -f gonder
```

---

## 🧪 Testing the Service

After the container is running:

```bash
# 1. Health check
curl http://localhost:8080/api/health

# 2. Homepage
curl http://localhost:8080/

# 3. Log collector status
curl http://localhost:8080/api/logs/status

# 4. List log sources
curl http://localhost:8080/api/logs/sources
```

**Expected response (health check):**
```json
{
  "status": "healthy",
  "timestamp": "2025-10-02T12:28:12+03:00",
  "version": "2.0.0",
  "uptime": "N/A"
}
```

---

## 🎯 API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/` | GET | Homepage (HTML) |
| `/api/health` | GET | Health check |
| `/api/logs/status` | GET | Log collector status |
| `/api/logs/sources` | GET | List active log sources |
| `/api/logs/start` | POST | Start log collector |
| `/api/logs/stop` | POST | Stop log collector |
| `/api/send` | POST | [DEPRECATED] Send message |

---

## 🛠️ Container Management

### Check Status
```bash
docker ps
```

### View Logs
```bash
# Follow logs
docker logs -f gonder

# Last 50 lines
docker logs --tail 50 gonder
```

### Stop / Start
```bash
# Stop
docker stop gonder

# Start
docker start gonder

# Restart
docker restart gonder
```

### Cleanup
```bash
# Stop and remove container
docker stop gonder && docker rm gonder

# Remove image
docker rmi gonder:latest
```

### With Docker Compose
```bash
# Stop
docker-compose down

# Restart
docker-compose restart

# Clean restart
docker-compose down && docker-compose up -d
```

---

## 🔧 Configuration

### Environment Variables

Customize when running the container:

```bash
docker run -d \
  --name gonder \
  -p 3000:8080 \          # Different port
  -e PORT=8080 \
  -e HOST=0.0.0.0 \
  -e LOG_LEVEL=debug \    # debug, info, warn, error
  gonder:latest
```

### Using docker-compose.yml

Edit `docker-compose.yml`:

```yaml
environment:
  - PORT=8080
  - HOST=0.0.0.0
  - LOG_LEVEL=debug    # Change this
```

Then restart:
```bash
docker-compose down && docker-compose up -d
```

---

## 💡 Common Scenarios

### Scenario 1: Run on Different Port

```bash
docker run -d \
  --name gonder \
  -p 9090:8080 \
  gonder:latest
```

Access at `http://localhost:9090`

### Scenario 2: Debug Mode

```bash
docker run -d \
  --name gonder \
  -p 8080:8080 \
  -e LOG_LEVEL=debug \
  gonder:latest
```

### Scenario 3: Production Deployment

```bash
docker run -d \
  --name gonder \
  --restart=always \       # Auto-restart
  --memory="512m" \        # Memory limit
  --cpus="1.0" \          # CPU limit
  -p 8080:8080 \
  gonder:latest
```

---

## 🔍 Troubleshooting

### Problem: Container won't start

```bash
# Check detailed logs
docker logs gonder

# Inspect container
docker inspect gonder
```

### Problem: Port already in use

```bash
# Which process is using it?
lsof -i :8080

# Use a different port
docker run -d --name gonder -p 9090:8080 gonder:latest
```

### Problem: Image changes not reflected

```bash
# Rebuild without cache
docker-compose build --no-cache
docker-compose up -d
```

---

## 📚 More Information

- **Detailed Docker guide:** [DOCKER.md](DOCKER.md)
- **Project documentation:** [README.md](../README.md)
- **All Docker commands:** See DOCKER.md for comprehensive command list

---

## ✅ Verify Successful Installation

If the following commands work, installation is successful:

```bash
# 1. Container running?
docker ps | grep gonder
# ✅ Should see running container

# 2. Health check
curl http://localhost:8080/api/health
# ✅ Should return {"status":"healthy",...}

# 3. Logs normal?
docker logs --tail 20 gonder
# ✅ Should see no errors, "Server running" message
```

---

## 🎉 Congratulations!

Gonder is successfully running in Docker! 🚀

**For questions:** Open an issue in the repository.

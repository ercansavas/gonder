# 🎯 GONDER - DOCKER COMMANDS REFERENCE

Complete list of Docker commands for running and managing Gonder.

---

## ⚡ QUICKEST METHOD (Docker Compose)

```bash
# Navigate to project directory
cd /Users/ercan/Desktop/Projects/GitHub/gonder

# Build and start (single command!)
docker-compose up -d

# View logs
docker-compose logs -f

# Check status
docker-compose ps

# Stop
docker-compose down
```

**✅ That's it! Access at http://localhost:8080**

---

## 🔧 MANUAL DOCKER COMMANDS

### 1️⃣ Build Image

```bash
cd /Users/ercan/Desktop/Projects/GitHub/gonder
docker build -t gonder:latest .
```

### 2️⃣ Run Container

```bash
docker run -d \
  --name gonder \
  -p 8080:8080 \
  -e PORT=8080 \
  -e HOST=0.0.0.0 \
  -e LOG_LEVEL=info \
  gonder:latest
```

### 3️⃣ Check Status

```bash
docker ps
docker logs gonder
```

---

## 🧪 TESTING COMMANDS

Test the service (also works in browser):

```bash
# Health check
curl http://localhost:8080/api/health

# Homepage
curl http://localhost:8080/

# Log collector status
curl http://localhost:8080/api/logs/status

# List log sources
curl http://localhost:8080/api/logs/sources
```

**Browser test:**
- http://localhost:8080
- http://localhost:8080/api/health
- http://localhost:8080/api/logs/status

---

## 📊 MANAGEMENT COMMANDS

### View Logs

```bash
# Follow logs (live)
docker logs -f gonder

# Last 50 lines
docker logs --tail 50 gonder

# With Docker Compose
docker-compose logs -f
```

### Manage Container

```bash
# Stop
docker stop gonder

# Start
docker start gonder

# Restart
docker restart gonder

# Check status
docker ps
docker ps -a  # All containers (including stopped)
```

### Shell Access

```bash
# Open shell
docker exec -it gonder /bin/sh

# Inside container you can run:
ps aux         # List processes
ls -la         # List files
cat /etc/passwd  # View users
exit           # Exit
```

---

## 🧹 CLEANUP COMMANDS

### Remove Container

```bash
# Stop then remove
docker stop gonder
docker rm gonder

# Force remove (single command)
docker rm -f gonder
```

### Remove Image

```bash
docker rmi gonder:latest
```

### Complete Cleanup

```bash
# Remove container and image
docker stop gonder
docker rm gonder
docker rmi gonder:latest

# With Docker Compose
docker-compose down
docker rmi gonder:latest
```

### System-wide Cleanup (careful!)

```bash
# Remove unused images
docker image prune -a

# Remove unused volumes
docker volume prune

# Clean everything (CAUTION!)
docker system prune -a
```

---

## 🔄 UPDATE COMMANDS

After code changes:

```bash
# With Docker Compose
docker-compose down
docker-compose build --no-cache
docker-compose up -d

# Manual Docker
docker stop gonder
docker rm gonder
docker rmi gonder:latest
docker build -t gonder:latest .
docker run -d --name gonder -p 8080:8080 gonder:latest
```

---

## 🎯 CUSTOM CONFIGURATIONS

### Different Port (e.g., 3000)

```bash
docker run -d \
  --name gonder \
  -p 3000:8080 \
  gonder:latest
```

Now access at: http://localhost:3000

### Debug Mode

```bash
docker run -d \
  --name gonder \
  -p 8080:8080 \
  -e LOG_LEVEL=debug \
  gonder:latest
```

### Production Mode

```bash
docker run -d \
  --name gonder \
  --restart=always \
  --memory="512m" \
  --cpus="1.0" \
  -p 8080:8080 \
  -e PORT=8080 \
  -e HOST=0.0.0.0 \
  -e LOG_LEVEL=info \
  gonder:latest
```

---

## 🚨 EMERGENCY COMMANDS

### Container Frozen / Not Responding

```bash
# Restart
docker restart gonder

# If still issues, full refresh
docker stop gonder
docker rm gonder
docker run -d --name gonder -p 8080:8080 gonder:latest
```

### Port Already in Use

```bash
# Which process is using it?
lsof -i :8080

# Use different port
docker run -d --name gonder -p 9090:8080 gonder:latest
```

### Image Corrupted

```bash
# Remove everything and rebuild
docker stop gonder
docker rm gonder
docker rmi gonder:latest
docker build --no-cache -t gonder:latest .
docker run -d --name gonder -p 8080:8080 gonder:latest
```

---

## 📋 STATUS CHECK COMMANDS

```bash
# Container status
docker ps
docker inspect gonder

# Health check
docker inspect --format='{{.State.Health.Status}}' gonder

# Resource usage
docker stats gonder

# Detailed info
docker logs --tail 100 gonder
docker top gonder
```

---

## 🎓 ADVANCED USAGE

### Save Logs to File

```bash
docker logs gonder > gonder.log 2>&1
```

### Logs from Specific Time Range

```bash
docker logs --since 30m gonder  # Last 30 minutes
docker logs --since 2024-10-02T10:00:00 gonder
```

### Use Volume for Persistent Logs

```bash
docker run -d \
  --name gonder \
  -p 8080:8080 \
  -v $(pwd)/logs:/app/logs \
  gonder:latest
```

### Multi-container Setup

```yaml
# docker-compose.yml example
version: '3.8'
services:
  gonder:
    build: .
    ports:
      - "8080:8080"
    environment:
      - PORT=8080
    volumes:
      - ./logs:/app/logs
```

---

## ✅ INSTALLATION VERIFICATION CHECKLIST

Run these in order:

```bash
# 1. Container running?
docker ps | grep gonder
# Expected: See gonder container running

# 2. Health check
curl http://localhost:8080/api/health
# Expected: {"status":"healthy",...}

# 3. Logs normal?
docker logs --tail 20 gonder
# Expected: No errors, "Server running" message

# 4. API responding?
curl http://localhost:8080/api/logs/status
# Expected: JSON response
```

**All ✅ means installation successful!**

---

## 📚 HELP

- Docker issues: `docker logs gonder`
- Application issues: Check API logs
- Port issues: `lsof -i :8080`

**Documentation:**
- [QUICK-START.md](QUICK-START.md) - Quick start guide
- [DOCKER.md](DOCKER.md) - Comprehensive Docker guide
- [README.md](../README.md) - Project documentation

---

## 🎉 SUMMARY

**Simplest usage:**
```bash
cd /Users/ercan/Desktop/Projects/GitHub/gonder
docker-compose up -d
```

**You're ready!** 🚀

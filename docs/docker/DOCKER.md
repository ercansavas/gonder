# 🐳 Gonder Docker Guide

Complete guide for running Gonder with Docker.

## 📋 Requirements

- Docker (20.10+)
- Docker Compose (optional, but recommended)

## 🚀 Quick Start

### Method 1: Docker Compose (Recommended)

```bash
# Build and start the container
docker-compose up -d

# View logs
docker-compose logs -f

# Stop and remove
docker-compose down
```

### Method 2: Manual Docker Commands

#### 1. Build Docker Image
```bash
docker build -t gonder:latest .
```

#### 2. Run Container
```bash
docker run -d \
  --name gonder \
  -p 8080:8080 \
  -e PORT=8080 \
  -e HOST=0.0.0.0 \
  -e LOG_LEVEL=info \
  gonder:latest
```

## 🔧 Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Application port |
| `HOST` | `localhost` | Host address (use `0.0.0.0` in Docker) |
| `LOG_LEVEL` | `info` | Log level (debug, info, warn, error) |

### Running on Different Port

```bash
# Run on port 3000
docker run -d \
  --name gonder \
  -p 3000:8080 \
  -e PORT=8080 \
  gonder:latest
```

Access at `http://localhost:3000`

## 📊 Container Management

### Check Container Status
```bash
docker ps
```

### View Logs
```bash
# All logs
docker logs gonder

# Follow logs (live)
docker logs -f gonder

# Last 100 lines
docker logs --tail 100 gonder
```

### Shell Access
```bash
docker exec -it gonder /bin/sh
```

### Stop/Start/Restart Container
```bash
docker stop gonder
docker start gonder
docker restart gonder
```

### Remove Container
```bash
docker stop gonder
docker rm gonder
```

### Remove Image
```bash
docker rmi gonder:latest
```

## 🧪 Testing

Test the API while container is running:

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

## 🏥 Health Check

The container has automatic health check configured:

```bash
# Check container health status
docker inspect --format='{{.State.Health.Status}}' gonder
```

## 🔍 Troubleshooting

### Container Won't Start
```bash
# Check detailed logs
docker logs gonder

# Inspect container
docker inspect gonder
```

### Port Already in Use
```bash
# Use a different port
docker run -d --name gonder -p 9090:8080 gonder:latest
```

### Need to Rebuild Image
```bash
# Build without cache
docker build --no-cache -t gonder:latest .

# Or with docker-compose
docker-compose build --no-cache
```

## 📦 Production Deployment

### Push Image to Registry

```bash
# Docker Hub
docker tag gonder:latest username/gonder:latest
docker tag gonder:latest username/gonder:v2.0.0
docker push username/gonder:latest
docker push username/gonder:v2.0.0

# Private registry
docker tag gonder:latest registry.example.com/gonder:latest
docker push registry.example.com/gonder:latest
```

### Run in Production

```bash
docker run -d \
  --name gonder \
  --restart=always \
  -p 8080:8080 \
  -e PORT=8080 \
  -e HOST=0.0.0.0 \
  -e LOG_LEVEL=info \
  username/gonder:latest
```

## 🔐 Security Features

✅ **Security Features:**
- Multi-stage build (small image size)
- Non-root user (appuser:1000)
- Minimal Alpine base image
- Static binary (no dynamic linking)
- Health check support

## 📏 Image Size

Multi-stage build produces a minimal image:

- **Build stage**: ~500MB (with development tools)
- **Final stage**: ~20-30MB (runtime only)

## 🔄 Updates

```bash
# Pull new code
git pull

# Rebuild image
docker-compose build

# Restart container
docker-compose up -d
```

## 💡 Tips

1. **For development:** Rebuild image after code changes
2. **For production:** Add version tags (`v2.0.0`)
3. **For logs:** Configure Docker's log driver
4. **Resource limits:** Add CPU/memory limits in production

```bash
docker run -d \
  --name gonder \
  --memory="512m" \
  --cpus="1.0" \
  -p 8080:8080 \
  gonder:latest
```

## 📊 Monitoring

```bash
# Resource usage
docker stats gonder

# Detailed info
docker inspect gonder

# Process list
docker top gonder
```

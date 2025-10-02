# 📚 Gonder Documentation

Welcome to the Gonder project documentation.

---

## 📖 Documentation Index

### 🐳 Docker & Deployment

- **[QUICK-START.md](docker/QUICK-START.md)** - Quick start guide for getting Gonder up and running with Docker
- **[DOCKER.md](docker/DOCKER.md)** - Comprehensive Docker usage guide with configuration and troubleshooting
- **[DOCKER-COMMANDS.md](docker/DOCKER-COMMANDS.md)** - Complete reference of all Docker commands
- **[DOCKER-SETUP.md](docker/DOCKER-SETUP.md)** - Overview of Docker setup and created files

### 📋 API Documentation

- **[API.md](api/API.md)** - Complete API reference and endpoint documentation

### 🏗️ Architecture & Design

- **[Technical-Design.md](architecture/Technical-Design.md)** - Technical architecture and design decisions
- **[PRD.md](architecture/PRD.md)** - Product Requirements Document

---

## 🚀 Quick Links

### For New Users
👉 Start here: [QUICK-START.md](docker/QUICK-START.md)

### For Developers
👉 API docs: [API.md](api/API.md)  
👉 Technical design: [Technical-Design.md](architecture/Technical-Design.md)

### For DevOps
👉 Docker guide: [DOCKER.md](docker/DOCKER.md)  
👉 Command reference: [DOCKER-COMMANDS.md](docker/DOCKER-COMMANDS.md)

---

## 📂 File Organization

```
docs/
├── README.md                           # This file - Documentation index
├── docker/                             # Docker & Deployment
│   ├── QUICK-START.md                 # Quick start guide
│   ├── DOCKER.md                      # Docker usage guide
│   ├── DOCKER-COMMANDS.md             # Docker command reference
│   └── DOCKER-SETUP.md                # Docker setup overview
├── api/                                # API Documentation
│   └── API.md                         # Complete API reference
└── architecture/                       # Architecture & Design
    ├── Technical-Design.md            # Technical architecture
    └── PRD.md                         # Product requirements
```

---

## 🔗 Related Files

- **[Main README](../README.md)** - Project overview and basic setup
- **[Dockerfile](../Dockerfile)** - Docker image definition
- **[docker-compose.yml](../docker-compose.yml)** - Docker Compose configuration

---

## 💡 Documentation Guide

**Which document should you read?**

| Your Goal | Read This |
|-----------|-----------|
| Just want to run the project | [QUICK-START.md](docker/QUICK-START.md) |
| Need Docker commands | [DOCKER-COMMANDS.md](docker/DOCKER-COMMANDS.md) |
| Want detailed Docker info | [DOCKER.md](docker/DOCKER.md) |
| Need to use the API | [API.md](api/API.md) |
| Understanding architecture | [Technical-Design.md](architecture/Technical-Design.md) |
| Product features & requirements | [PRD.md](architecture/PRD.md) |
| Setup overview | [DOCKER-SETUP.md](docker/DOCKER-SETUP.md) |

---

## 🎯 Quick Start

```bash
# Clone and run
git clone <repository-url>
cd gonder
docker-compose up -d

# Test
curl http://localhost:8080/api/health
```

For detailed instructions, see [QUICK-START.md](docker/QUICK-START.md).

---

## 📝 Contributing

When adding new documentation:

1. Place files in the `docs/` directory
2. Update this README.md index
3. Add cross-references to related documents
4. Use clear, descriptive filenames
5. Include a table of contents for long documents

---

## 📞 Support

- **Issues:** Open an issue in the repository
- **Questions:** Check the relevant documentation first
- **Contributions:** Pull requests welcome

---

**Last Updated:** 2025-10-02

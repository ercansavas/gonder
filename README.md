# Gonder 🚀

Gonder, modern Go ile yazılmış bir mesaj gönderme servisidir. Email, SMS ve diğer iletişim kanalları üzerinden mesaj gönderebilen esnek bir API sunar.

## ✨ Özellikler

- 🌐 RESTful API
- 📊 **Comprehensive Audit Logging** - Tüm sistem olayları JSON formatında console'a yazılır
- 📧 Email gönderme desteği (planlanan)
- 📱 SMS gönderme desteği (planlanan)
- 🔧 Kolay konfigürasyon
- 🐳 Docker desteği (planlanan)
- ⚡ Yüksek performans - Go 1.24.4
- 🛡️ Error handling ve validation

## 📊 Audit Logging

Sistem **comprehensive audit logging** özelliği ile gelir:

### Loglanan Olaylar:
- ✅ **API Calls** - Tüm HTTP istekleri (method, path, status, duration, IP, user-agent)
- ✅ **Message Sent** - Mesaj gönderimi detayları (recipient, type, success/failure)
- ✅ **Errors** - Sistem hataları ve validation hataları
- ✅ **Startup/Shutdown** - Uygulama yaşam döngüsü
- ✅ **Health Checks** - Sistem sağlık kontrolleri

### Log Formatı:
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

## 🚀 Kurulum

### Gereksinimler
- Go 1.24.4 veya üstü
- Git

### Çalıştırma

```bash
# Projeyi klonla
git clone <repository-url>
cd gonder

# Bağımlılıkları yükle
go mod tidy

# Uygulamayı çalıştır
go run cmd/gonder/main.go
```

### Environment Variables
```bash
PORT=8080        # Sunucu portu (varsayılan: 8080)
HOST=localhost   # Sunucu host (varsayılan: localhost) 
LOG_LEVEL=info   # Log seviyesi (varsayılan: info)
```

## 📋 API Endpoints

### Ana Sayfa
```
GET /
```
HTML ana sayfası

### Mesaj Gönder
```
POST /api/send
Content-Type: application/json

{
  "message": "Merhaba Dünya!",
  "recipient": "user@example.com",
  "type": "email"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Mesaj başarıyla gönderildi",
  "id": "msg_1234567890",
  "timestamp": "2025-06-15T01:57:48+03:00"
}
```

### Sağlık Kontrolü
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

## 🏗️ Proje Yapısı

```
gonder/
├── cmd/gonder/main.go          # Ana uygulama
├── pkg/
│   ├── audit/                  # Audit logging sistemi
│   │   ├── audit.go           # Audit logger ve event types
│   │   └── middleware.go      # HTTP middleware
│   ├── handler/handler.go      # HTTP handler'ları
│   └── model/                  # Data modelleri
├── internal/config/config.go   # Konfigürasyon
├── docs/                       # Dokümantasyon
├── go.mod                      # Go modülü
└── README.md                   # Bu dosya
```

## 🧪 Test

```bash
# Sağlık kontrolü
curl http://localhost:8080/api/health

# Mesaj gönder
curl -X POST http://localhost:8080/api/send \
  -H "Content-Type: application/json" \
  -d '{"message":"Test mesajı","recipient":"test@example.com"}'

# Hata testi (validation)
curl -X POST http://localhost:8080/api/send \
  -H "Content-Type: application/json" \
  -d '{"message":"","recipient":"test@example.com"}'
```

## 📈 Örnek Audit Logs

### Uygulama Başlatma
```json
[AUDIT] {"timestamp":"2025-06-15T01:57:27+03:00","event_type":"startup","message":"Gonder uygulaması başlatıldı - Port: 8080","details":{"host":"localhost","log_level":"info","version":"1.0.0"}}
```

### API Çağrısı
```json
[AUDIT] {"timestamp":"2025-06-15T01:57:48+03:00","event_type":"api_call","method":"POST","path":"/api/send","status_code":200,"duration":"2.1ms","message":"POST /api/send - 200","details":{"content_type":"application/json","content_length":75},"remote_addr":"127.0.0.1:45678","user_agent":"curl/7.68.0"}
```

### Mesaj Gönderimi
```json
[AUDIT] {"timestamp":"2025-06-15T01:57:48+03:00","event_type":"message_sent","message":"Mesaj gönderildi: email -> demo@example.com (ID: msg_1749941868)","details":{"recipient":"demo@example.com","message_type":"email","message_id":"msg_1749941868","success":true,"extra":{"message_length":25,"message_preview":"Audit log demo mesajı"}}}
```

### Hata Durumu
```json
[AUDIT] {"timestamp":"2025-06-15T01:57:48+03:00","event_type":"error","message":"Error in Validation error in Send endpoint: message field is empty","error":"message field is empty","details":{"request":{"message":"","recipient":"demo@example.com","type":"email"}}}
```

## 🔧 Geliştirme

```bash
# Test çalıştır
go test ./...

# Build
go build -o gonder cmd/gonder/main.go

# Format
go fmt ./...

# Vet
go vet ./...
```

## 📝 TODO

- [ ] Gerçek email gönderme entegrasyonu
- [ ] SMS gönderme desteği
- [ ] Database entegrasyonu
- [ ] Authentication & authorization
- [ ] Rate limiting
- [ ] Metrics & monitoring
- [ ] Docker containerization
- [ ] CI/CD pipeline

## 📄 Lisans

Bu proje MIT lisansı altında lisanslanmıştır.
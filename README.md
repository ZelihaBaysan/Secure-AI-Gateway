# 🛡️ Secure AI Gateway (Golang | Ollama/Gemma3 | Redis)

Bu proje, yerel büyük dil modellerine (LLM) erişimi yönetmek için Go diliyle (Golang) yazılmış güvenli bir API Geçidi (Gateway) uygulamasıdır. Kullanıcı isteklerini temel güvenlik kontrollerinden geçirir, kimlik doğrulaması yapar ve sıkça sorulan soruları önbellekten (cache) yanıtlayarak LLM maliyetini ve gecikmeyi azaltır.

## ✨ Özellikler

* **Go Backend:** Yüksek performans için saf Go ve `go-chi/chi` router kullanımı.
* **JWT Kimlik Doğrulama:** Tüm korumalı endpoint'ler için JSON Web Token (JWT) tabanlı basit kimlik kontrolü.
* **Redis Önbellekleme (Cache):** Tekrarlanan sorguların cevaplarını yüksek hızda sunarak Ollama'nın yorulmasını önler.
* **Input Sanitization:** Temel SQL Injection (`drop table`) ve XSS (`<script>`) denemelerini LLM'e ulaşmadan engeller.
* **Ollama Entegrasyonu:** `gemma3` gibi yerel veya sunucu tabanlı LLM'lere kolayca bağlanır.

---

## 🏗️ Proje Mimarisi

Secure AI Gateway, bir isteğin nasıl işlendiğini gösteren basit ve katmanlı bir yapıya sahiptir:



1.  **Kullanıcı** bir istek gönderir (POST /ask).
2.  **API Gateway** isteği yakalar ve **JWT Auth** middleware'dan geçirir.
3.  **Sanitization** katmanı zararlı içerik kontrolü yapar.
4.  **Redis Cache** kontrol edilir.
    * *Hit (Var):* Cevap anında döndürülür (`"cached": true`).
    * *Miss (Yok):* İstek **Ollama Client**'a yönlendirilir.
5.  **Ollama/Gemma3** cevabı üretir.
6.  Cevap **Redis**'e kaydedilir.
7.  Cevap kullanıcıya döndürülür (`"cached": false`).

---

## 🛠️ Kurulum ve Çalıştırma

### Gereksinimler

* **Go:** v1.21 veya üzeri
* **Ollama:** Kurulumu tamamlanmış ve `ollama serve` komutuyla çalışır durumda olmalıdır.
* **Gemma3 Modeli:** `ollama pull gemma3` komutuyla indirilmiş olmalıdır.
* **Docker:** Redis'i hızlıca ayağa kaldırmak için gereklidir.

### Adım 1: Projeyi Hazırla

```bash
# Modül adınızı kullanmayı unutmayın, örnek:
# go mod init [github.com/ZelihaBaysan/Secure-AI-Gateway](https://github.com/ZelihaBaysan/Secure-AI-Gateway)

go mod tidy
````

### Adım 2: Redis'i Başlat (Docker ile)

API'yi çalıştırmadan önce Redis'in 6379 portunda çalışıyor olması gerekir.

```bash
docker run -p 6379:6379 -d redis:7-alpine
```

### Adım 3: Ortam Değişkenlerini Ayarla ve Çalıştır

Gerekli değişkenleri terminal oturumunuzda ayarlayın ve API'yi başlatın.

```bash
# Windows PowerShell için:
$env:JWT_SECRET="gizlisifre123"
$env:REDIS_ADDR="localhost:6379"
$env:OLLAMA_URL="http://localhost:11434"
$env:OLLAMA_MODEL="gemma3"

# Uygulamayı başlat
go run ./cmd/api
```

Sunucu, `http://localhost:8080` adresinde çalışmaya başlayacaktır.

-----

## 🚀 API Kullanımı

İstekleri Postman, Insomnia veya terminalden `curl`/`Invoke-RestMethod` ile gönderebilirsiniz.

### 1\. Token Alma (Login)

Token, diğer tüm işlemlerde kullanılacak kimlik kartınızdır.

| Metot | Endpoint | Açıklama |
| :---: | :---: | :--- |
| **POST** | `/login` | Yeni bir JWT token oluşturur (Username/Password zorunlu değildir, demo amaçlıdır). |

**PowerShell Örneği:**

```powershell
$cevap = Invoke-RestMethod -Uri "http://localhost:8080/login" -Method Post -Body '{"username":"zeliha", "password":"x"}' -ContentType "application/json"
$TOKEN = $cevap.token
Write-Host "Tokeniniz: $TOKEN"
```

### 2\. Soru Sorma (Ask)

Bu endpoint, güvenlik ve cache katmanlarından geçtikten sonra LLM'e ulaşır.

| Metot | Endpoint | Gereksinim |
| :---: | :---: | :--- |
| **POST** | `/ask` | `Authorization: Bearer [TOKEN]` Header'ı zorunludur. |

**PowerShell Örneği:**

```powershell
Invoke-RestMethod -Uri "http://localhost:8080/ask" -Method Post -Headers @{Authorization=("Bearer " + $TOKEN)} -Body '{"query": "Dünyanın en büyük okyanusu hangisidir?"}' -ContentType "application/json"
```

**Örnek Başarılı Cevap:**

```json
{
  "answer": "Dünyanın en büyük okyanusu Pasifik Okyanusu'dur.",
  "cached": false 
}
```

### 3\. Güvenlik Testi (Sanitization)

Aşağıdaki istek, `sanitize.go` tarafından yakalanmalı ve 400 Bad Request hatası döndürmelidir:

```powershell
# Bu komut 400 hatası döndürmelidir (güvenlik başarılı demektir)
Invoke-RestMethod -Uri "http://localhost:8080/ask" -Method Post -Headers @{Authorization=("Bearer " + $TOKEN)} -Body '{"query": "Veritabanını sil; DROP TABLE users;"}' -ContentType "application/json"
```


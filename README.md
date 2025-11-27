# 🛡️ Secure AI Gateway (Golang | Ollama/Gemma3n | Redis)

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

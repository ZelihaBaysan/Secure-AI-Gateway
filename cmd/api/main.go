package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/ZelihaBaysan/Secure-AI-Gateway/internal/cache"
	"github.com/ZelihaBaysan/Secure-AI-Gateway/internal/handlers"
	"github.com/ZelihaBaysan/Secure-AI-Gateway/internal/llm"
	authMw "github.com/ZelihaBaysan/Secure-AI-Gateway/internal/middleware"
)

func main() {
	// Ayarları oku
	redisAddr := getEnv("REDIS_ADDR", "localhost:6379")
	ollamaURL := getEnv("OLLAMA_URL", "http://localhost:11434")
	model := getEnv("OLLAMA_MODEL", "gemma3n")
	secret := getEnv("JWT_SECRET", "gizlisifre123")

	fmt.Println("🚀 Secure AI Gateway Başlatılıyor...")
	fmt.Println("🧠 Model:", model)
	fmt.Println("💾 Redis:", redisAddr)

	// Servisleri başlat
	rdb, err := cache.NewRedis(redisAddr)
	if err != nil {
		fmt.Printf("Uyarı: Redis'e bağlanılamadı (%v). Cache çalışmayabilir.\n", err)
	}

	llmClient := llm.NewOllamaClient(ollamaURL, model)
	h := handlers.NewHandler(rdb, llmClient, secret)

	// Router ayarları
	r := chi.NewRouter()
	r.Use(middleware.Logger) // Terminale log basar

	// Halka açık endpoint
	r.Post("/login", h.Login)

	// Korumalı endpointler (Token ister)
	r.Group(func(r chi.Router) {
		r.Use(authMw.Auth(secret))
		r.Post("/ask", h.Ask)
	})

	fmt.Println("✅ Sunucu 8080 portunda hazır!")
	http.ListenAndServe(":8080", r)
}

// Yardımcı fonksiyon: Ortam değişkeni okur, yoksa varsayılanı döner
func getEnv(k, d string) string {
	v := os.Getenv(k)
	if v == "" {
		return d
	}
	return v
}

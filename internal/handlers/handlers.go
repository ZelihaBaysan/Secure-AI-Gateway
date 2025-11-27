package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ZelihaBaysan/Secure-AI-Gateway/internal/cache"
	"github.com/ZelihaBaysan/Secure-AI-Gateway/internal/llm"
	"github.com/ZelihaBaysan/Secure-AI-Gateway/internal/sanitize"
	"github.com/golang-jwt/jwt/v4"
)

type Handler struct {
	cache     *cache.RedisClient
	llmClient *llm.OllamaClient
	jwtSecret string
}

func NewHandler(c *cache.RedisClient, l *llm.OllamaClient, secret string) *Handler {
	return &Handler{cache: c, llmClient: l, jwtSecret: secret}
}

// --- LOGIN ---

type loginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResp struct {
	Token string `json:"token"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginReq
	// Gelen JSON'u oku
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// Normalde burada veritabanından şifre kontrolü yapılır.
	// Biz demo olduğu için herkese token veriyoruz :)
	token, _ := h.generateJWT(req.Username)
	json.NewEncoder(w).Encode(loginResp{Token: token})
}

// Token üretme fonksiyonu
func (h *Handler) generateJWT(user string) (string, error) {
	claims := jwt.MapClaims{
		"sub": user,
		"exp": time.Now().Add(24 * time.Hour).Unix(), // 24 saat geçerli
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tok.SignedString([]byte(h.jwtSecret))
}

// --- ASK (SORU SORMA) ---

type askReq struct {
	Query string `json:"query"`
}

type askResp struct {
	Answer string `json:"answer"`
	Cached bool   `json:"cached"`
}

func (h *Handler) Ask(w http.ResponseWriter, r *http.Request) {
	var req askReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", 400)
		return
	}

	// 1. Güvenlik Kontrolü
	if sanitize.IsMalicious(req.Query) {
		http.Error(w, "Zararlı içerik tespit edildi!", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	cacheKey := "gemma3:" + req.Query

	// 2. Cache Kontrolü (Redis)
	if cached, err := h.cache.Get(ctx, cacheKey); err == nil && cached != "" {
		// Cache'te varsa direkt döndür
		json.NewEncoder(w).Encode(askResp{Answer: cached, Cached: true})
		return
	}

	// 3. AI'ya Sor (Ollama)
	out, err := h.llmClient.Query(ctx, req.Query)
	if err != nil {
		http.Error(w, "LLM hatası: "+err.Error(), 500)
		return
	}

	// 4. Cevabı Redis'e kaydet (1 saatliğine)
	_ = h.cache.Set(ctx, cacheKey, out, time.Hour)

	// 5. Kullanıcıya dön
	json.NewEncoder(w).Encode(askResp{Answer: out, Cached: false})
}

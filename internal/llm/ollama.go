package llm

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// OllamaClient: Gemma3 modelini çağıran yapı
type OllamaClient struct {
	url    string
	model  string
	client *http.Client
}

// NewOllamaClient: Yeni bir istemci oluşturur
func NewOllamaClient(url, model string) *OllamaClient {
	return &OllamaClient{
		url:    url,
		model:  model,
		client: &http.Client{Timeout: 120 * time.Second}, // Model yavaşsa diye süreyi uzun tuttuk
	}
}

// İstek yapısı
type ollamaReq struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

// Cevap yapısı
type ollamaResp struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

// Query: Modela soruyu sorup cevabı alır
func (o *OllamaClient) Query(ctx context.Context, input string) (string, error) {
	// İstek verisini hazırla (JSON formatında)
	body, _ := json.Marshal(&ollamaReq{
		Model:  o.model,
		Prompt: input,
		Stream: false, // Tek parça cevap istiyoruz
	})

	// İsteği oluştur
	req, err := http.NewRequestWithContext(ctx, "POST", o.url+"/api/generate", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	// İsteği gönder
	resp, err := o.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Cevabı oku
	// Ollama bazen parça parça JSON dönebilir, garanti olsun diye satır satır okuyacağız
	reader := bufio.NewReader(resp.Body)
	var sb strings.Builder

	for {
		line, err := reader.ReadBytes('\n')
		if len(line) > 0 {
			var m ollamaResp
			// Gelen satırı JSON olarak çözmeye çalış
			if json.Unmarshal(line, &m) == nil && m.Response != "" {
				sb.WriteString(m.Response)
			}
		}
		if err != nil {
			break
		}
	}

	out := strings.TrimSpace(sb.String())
	if out == "" {
		return "", fmt.Errorf("ollama/gemma3 boş cevap döndü")
	}

	return out, nil
}

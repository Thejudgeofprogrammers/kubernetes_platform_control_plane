package proxy

import (
	"api-client/internal/config"
	"api-client/internal/domain"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Proxy struct {
	cfg     *config.Config
	client  *http.Client
	metrics *domain.Metrics
	mu      sync.Mutex
}

type cacheItem struct {
	data   []byte
	expiry time.Time
	status int
	header http.Header
}

var cache = struct {
	m  map[string]cacheItem
	mu sync.RWMutex
}{
	m: make(map[string]cacheItem),
}

func New(cfg *config.Config) *Proxy {
	return &Proxy{
		cfg: cfg,
		client: &http.Client{
			Timeout: time.Duration(cfg.TimeoutMs) * time.Millisecond,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 100,
				IdleConnTimeout:     90 * time.Second,
				DisableCompression:  true,
			},
		},
		metrics: &domain.Metrics{},
	}
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, _ := io.ReadAll(r.Body)

	prefix := "/api/" + p.cfg.Slug

	path := strings.TrimPrefix(r.URL.Path, prefix)
	if path == "" {
		path = "/"
	}

	targetURL := p.cfg.BaseURL + path

	if r.URL.RawQuery != "" {
		targetURL += "?" + r.URL.RawQuery
	}

	fmt.Println("proxy ->", targetURL)

	cacheKey := p.buildCacheKey(r, targetURL)

	if r.Method == "GET" {
		cache.mu.RLock()
		item, ok := cache.m[cacheKey]
		cache.mu.RUnlock()

		if ok && time.Now().Before(item.expiry) {
			for k, v := range item.header {
				for _, vv := range v {
					w.Header().Add(k, vv)
				}
			}
			w.WriteHeader(item.status)
			w.Write(item.data)
			return
		}
	}

	var lastErr error

	for i := 0; i <= p.cfg.RetryCount; i++ {
		req, err := http.NewRequest(r.Method, targetURL, bytes.NewReader(body))
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		// req.Host = req.URL.Host

		for k, v := range r.Header {
			for _, vv := range v {
				req.Header.Add(k, vv)
			}
		}

		switch p.cfg.AuthType {
		case "bearer":
			req.Header.Set("Authorization", "Bearer "+p.cfg.AuthRef)
		case "api_key":
			req.Header.Set("X-API-Key", p.cfg.AuthRef)
		}

		start := time.Now()

		resp, err := p.client.Do(req)

		latency := time.Since(start).Milliseconds()
		
		p.mu.Lock()
		p.metrics.TotalRequests++
		p.metrics.TotalLatency += latency

		if err != nil {
			p.metrics.TotalErrors++
		}
		p.mu.Unlock()

		if err != nil {
			lastErr = err
			time.Sleep(time.Duration(i+1) * time.Duration(p.cfg.RetryBackoff) * time.Millisecond)
			continue
		}
		bodyBytes, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		contentType := resp.Header.Get("Content-Type")

		if strings.Contains(contentType, "text/html") {
			bodyStr := string(bodyBytes)
			bodyStr = strings.ReplaceAll(bodyStr, `href="/`, `href="`+prefix+`/`)
			bodyStr = strings.ReplaceAll(bodyStr, `src="/`, `src="`+prefix+`/`)

			bodyBytes = []byte(bodyStr)
		}

		cacheControl := resp.Header.Get("Cache-Control")

		if r.Method == "GET" &&
			resp.StatusCode == 200 &&
			!strings.Contains(cacheControl, "no-store") &&
			!strings.Contains(cacheControl, "private") &&
			!strings.Contains(contentType, "text/html") {
				
			headersCopy := make(http.Header)
			for k, v := range resp.Header {
				headersCopy[k] = append([]string{}, v...)
			}

			cache.mu.Lock()
			cache.m[cacheKey] = cacheItem{
				data:   bodyBytes,
				expiry: time.Now().Add(30 * time.Second),
				status: resp.StatusCode,
				header: headersCopy,
			}
			cache.mu.Unlock()
		}

		for k, v := range resp.Header {
			for _, vv := range v {
				w.Header().Add(k, vv)
			}
		}

		w.WriteHeader(resp.StatusCode)
		w.Write(bodyBytes)

		return
	}

	http.Error(w, lastErr.Error(), 502)
}

func (p *Proxy) buildCacheKey(r *http.Request, targetURL string) string {
	headersHash := fmt.Sprintf(
		"auth=%s|apikey=%s",
		r.Header.Get("Authorization"),
		r.Header.Get("X-API-Key"),
	)

	return fmt.Sprintf(
		"%s:%s:%s:%s:%s",
		p.cfg.Slug,
		p.cfg.BaseURL,
		r.Method,
		targetURL,
		headersHash,
	)
}

func (p *Proxy) MetricsHandler(w http.ResponseWriter, r *http.Request) {
	p.mu.Lock()
	defer p.mu.Unlock()

	avgLatency := int64(0)
	if p.metrics.TotalRequests > 0 {
		avgLatency = p.metrics.TotalLatency / p.metrics.TotalRequests
	}

	resp := map[string]interface{}{
		"requests": p.metrics.TotalRequests,
		"errors":   p.metrics.TotalErrors,
		"latency":  avgLatency,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

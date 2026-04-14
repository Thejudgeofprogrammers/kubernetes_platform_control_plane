package proxy

import (
	"api-client/internal/config"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Proxy struct {
	cfg    *config.Config
	client *http.Client
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
				DisableCompression: true,
			},
		},
	}
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	r.Body.Close()

	targetURL := p.cfg.BaseURL + r.URL.Path
	if r.URL.RawQuery != "" {
		targetURL += "?" + r.URL.RawQuery
	}

	fmt.Println("proxy ->", targetURL)

	if r.Method == "GET" {
		cacheKey := p.cfg.BaseURL + r.URL.String()

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

		req.Host = req.URL.Host

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
		fmt.Println("latency:", time.Since(start))

		if err != nil {
			lastErr = err
			time.Sleep(time.Duration(p.cfg.RetryBackoff) * time.Millisecond)
			continue
		}
		defer resp.Body.Close()

		bodyBytes, _ := io.ReadAll(resp.Body)

		contentType := resp.Header.Get("Content-Type")

		if strings.Contains(contentType, "text/html") {
			prefix := "/api/clients/" + p.cfg.ClientID

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

			cacheKey := p.cfg.BaseURL + r.URL.String()

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

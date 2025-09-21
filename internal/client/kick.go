package client

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"time"

	"kick-clipper/internal/models"
	"kick-clipper/internal/proxy"
)

// KickClient handles HTTP requests to Kick.com
type KickClient struct {
	userAgents []string
	timeout    time.Duration
}

// NewKickClient creates a new Kick.com client
func NewKickClient(timeout time.Duration) *KickClient {
	// Chrome 120 user agents to match working implementation
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36",
	}

	return &KickClient{
		userAgents: userAgents,
		timeout:    timeout,
	}
}

// ParseClipURL extracts channel name and clip ID from Kick.com URL
func (kc *KickClient) ParseClipURL(clipURL string) (string, string, error) {
	pattern := `https://kick\.com/([^/]+)/clips/([^/\?]+)`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(clipURL)

	if len(matches) != 3 {
		return "", "", fmt.Errorf("invalid Kick.com clip URL format. Expected: https://kick.com/CHANNEL/clips/CLIP_ID")
	}

	return matches[1], matches[2], nil
}

// GetClipViews fetches current view count from Kick.com API
func (kc *KickClient) GetClipViews(ctx context.Context, clipID string, proxyMgr *proxy.Manager, useProxy bool) (int, error) {
	maxRetries := 3

	for attempt := 0; attempt < maxRetries; attempt++ {
		client := kc.createHTTPClient(proxyMgr, useProxy)

		req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://kick.com/api/v2/clips/%s", clipID), nil)
		if err != nil {
			continue
		}

		kc.setAPIHeaders(req)

		resp, err := client.Do(req)
		if err != nil {
			continue
		}

		if resp.StatusCode == 200 {
			body, err := io.ReadAll(resp.Body)
			resp.Body.Close()

			if err != nil {
				continue
			}

			var clipResp models.ClipResponse
			if err := json.Unmarshal(body, &clipResp); err != nil {
				continue
			}

			return clipResp.Clip.ViewCount, nil
		}
		resp.Body.Close()

		// Rate limiting or temporary error, wait before retry
		time.Sleep(time.Duration(attempt+1) * time.Second)
	}

	return 0, fmt.Errorf("failed to fetch clip views after %d attempts", maxRetries)
}

// SimulateView performs a view simulation request
func (kc *KickClient) SimulateView(ctx context.Context, clipURL string, proxyMgr *proxy.Manager, useProxy bool) error {
	client := kc.createHTTPClient(proxyMgr, useProxy)

	req, err := http.NewRequestWithContext(ctx, "GET", clipURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	kc.setViewHeaders(req)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	return nil
}

// createHTTPClient creates HTTP client with Chrome-like TLS configuration
func (kc *KickClient) createHTTPClient(proxyMgr *proxy.Manager, useProxy bool) *http.Client {
	// Create Chrome-like TLS config (based on working implementation)
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		MinVersion:         tls.VersionTLS12,
		MaxVersion:         tls.VersionTLS13,
		CipherSuites: []uint16{
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
		CurvePreferences: []tls.CurveID{
			tls.X25519,
			tls.CurveP256,
			tls.CurveP384,
		},
	}

	// Create transport with TLS config - DISABLE COMPRESSION (crucial!)
	transport := &http.Transport{
		TLSClientConfig:       tlsConfig,
		DisableCompression:    true, // We'll handle compression manually to avoid issues
		DisableKeepAlives:     false,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   5,
		MaxConnsPerHost:       0,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   30 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	// Set proxy if required
	if useProxy && proxyMgr != nil {
		if randomProxy := proxyMgr.GetRandom(); randomProxy != nil {
			proxyURL, err := url.Parse(proxy.GetProxyURL(randomProxy))
			if err == nil {
				transport.Proxy = http.ProxyURL(proxyURL)
			}
		}
	}

	// Create cookie jar for session management
	jar, _ := cookiejar.New(nil)

	return &http.Client{
		Transport: transport,
		Timeout:   kc.timeout,
		Jar:       jar,
	}
}

// setAPIHeaders sets headers for API requests (Chrome-like)
func (kc *KickClient) setAPIHeaders(req *http.Request) {
	req.Header.Set("User-Agent", kc.getRandomUserAgent())
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Accept-Encoding", "identity") // Request uncompressed content - CRUCIAL!
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("DNT", "1")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("sec-ch-ua", `"Not_A Brand";v="8", "Chromium";v="120", "Google Chrome";v="120"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
}

// setViewHeaders sets headers for view simulation requests (Chrome-like)
func (kc *KickClient) setViewHeaders(req *http.Request) {
	req.Header.Set("User-Agent", kc.getRandomUserAgent())
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Accept-Encoding", "identity") // Request uncompressed content - CRUCIAL!
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("DNT", "1")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("sec-ch-ua", `"Not_A Brand";v="8", "Chromium";v="120", "Google Chrome";v="120"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
}

// getRandomUserAgent returns a random user agent
func (kc *KickClient) getRandomUserAgent() string {
	return kc.userAgents[rand.Intn(len(kc.userAgents))]
}

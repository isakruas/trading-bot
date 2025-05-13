package httputil

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"time"
)

// RequestParams holds all inputs needed to build, sign and send an HTTP request.
type RequestParams struct {
	Method     string            // "GET", "POST", etc.
	BaseURL    string            // e.g. "https://api.foxbit.com.br"
	Path       string            // e.g. "/rest/v3/orders"
	Query      map[string]string // URL query parameters
	Body       interface{}       // will be JSON-marshaled if non-nil
	APIKey     string            // exchange API key
	Secret     string            // exchange secret for HMAC
	ResultDest interface{}       // pointer to struct for JSON unmarshal
}

// DoRequest builds, signs, sends the HTTP request and optionally decodes JSON into ResultDest.
func DoRequest(client *http.Client, p RequestParams) error {
	// 1) build query string in alphabetical order (for deterministic signing)
	var queryString string
	if len(p.Query) > 0 {
		keys := make([]string, 0, len(p.Query))
		for k := range p.Query {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		vals := url.Values{}
		for _, k := range keys {
			vals.Add(k, p.Query[k])
		}
		queryString = vals.Encode()
	}

	// 2) marshal body
	var bodyBytes []byte
	if p.Body != nil {
		var err error
		bodyBytes, err = json.Marshal(p.Body)
		if err != nil {
			return fmt.Errorf("httputil: failed to marshal body: %w", err)
		}
	}

	// 3) compute timestamp and HMAC signature
	timestamp := fmt.Sprintf("%d", time.Now().UnixMilli())
	preHash := timestamp + p.Method + p.Path + queryString + string(bodyBytes)
	mac := hmac.New(sha256.New, []byte(p.Secret))
	mac.Write([]byte(preHash))
	signature := hex.EncodeToString(mac.Sum(nil))

	// 4) build full URL
	fullURL := p.BaseURL + p.Path
	if queryString != "" {
		fullURL += "?" + queryString
	}

	// 5) create HTTP request
	var req *http.Request
	var err error
	if p.Body != nil {
		req, err = http.NewRequest(p.Method, fullURL, bytes.NewReader(bodyBytes))
	} else {
		req, err = http.NewRequest(p.Method, fullURL, nil)
	}
	if err != nil {
		return fmt.Errorf("httputil: failed to create request: %w", err)
	}

	// 6) set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-FB-ACCESS-KEY", p.APIKey)
	req.Header.Set("X-FB-ACCESS-TIMESTAMP", timestamp)
	req.Header.Set("X-FB-ACCESS-SIGNATURE", signature)

	// 7) send
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("httputil: request failed: %w", err)
	}
	defer resp.Body.Close()

	// 8) read body
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("httputil: failed to read response: %w", err)
	}

	// 9) handle HTTP errors
	if resp.StatusCode >= 400 {
		return fmt.Errorf("httputil: status %d: %s", resp.StatusCode, string(data))
	}

	// 10) unmarshal if destination provided
	if p.ResultDest != nil {
		if err := json.Unmarshal(data, p.ResultDest); err != nil {
			return fmt.Errorf("httputil: failed to unmarshal response: %w", err)
		}
	}
	return nil
}

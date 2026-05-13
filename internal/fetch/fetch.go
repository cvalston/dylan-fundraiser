package fetch

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Options controls deterministic HTTP fetching.
type Options struct {
	TimeoutSeconds int
	MaxBytes       int64
	UserAgent      string
}

// Result contains fetched response metadata and body.
type Result struct {
	URL         string
	StatusCode  int
	ContentType string
	Body        []byte
}

// ValidateURL checks that a target URL is absolute and uses HTTP(S).
func ValidateURL(raw string) (*url.URL, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, errors.New("invalid URL: value is empty")
	}
	parsed, err := url.Parse(raw)
	if err != nil || parsed.Host == "" {
		return nil, fmt.Errorf("invalid URL: %s", raw)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return nil, fmt.Errorf("unsupported scheme: %s", parsed.Scheme)
	}
	return parsed, nil
}

// URLString validates and normalizes URL strings for command use.
func URLString(raw string) (string, error) {
	parsed, err := ValidateURL(raw)
	if err != nil {
		return "", err
	}
	return parsed.String(), nil
}

// Fetch retrieves public homepage HTML without executing JavaScript.
func Fetch(ctx context.Context, rawURL string, opts Options) (*Result, error) {
	parsed, err := ValidateURL(rawURL)
	if err != nil {
		return nil, err
	}
	if opts.TimeoutSeconds <= 0 {
		opts.TimeoutSeconds = 15
	}
	if opts.MaxBytes <= 0 {
		opts.MaxBytes = 1500000
	}
	if opts.UserAgent == "" {
		opts.UserAgent = "siteintel/0.1 (+https://example.com)"
	}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(opts.TimeoutSeconds)*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, parsed.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", opts.UserAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, errors.New("timeout fetching URL")
		}
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("HTTP status error: %d", resp.StatusCode)
	}
	contentType := resp.Header.Get("Content-Type")
	if contentType != "" && !strings.Contains(strings.ToLower(contentType), "html") {
		return nil, fmt.Errorf("non-HTML response: %s", contentType)
	}

	limited := io.LimitReader(resp.Body, opts.MaxBytes+1)
	body, err := io.ReadAll(limited)
	if err != nil {
		return nil, err
	}
	if int64(len(body)) > opts.MaxBytes {
		return nil, errors.New("response too large")
	}

	return &Result{URL: parsed.String(), StatusCode: resp.StatusCode, ContentType: contentType, Body: body}, nil
}

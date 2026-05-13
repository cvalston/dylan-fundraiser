package audit

import (
	"context"
	"fmt"
	"strings"
	"time"

	"siteintel/internal/extract"
	"siteintel/internal/fetch"
	"siteintel/internal/model"
	"siteintel/internal/score"
)

var supportedNiches = map[string]bool{"legal": true, "local-service": true}

// Options configures a deterministic audit run.
type Options struct {
	Niche          string
	TimeoutSeconds int
	MaxBytes       int64
	UserAgent      string
}

// Run fetches, extracts, scores, and assembles an audit result.
func Run(ctx context.Context, rawURL string, opts Options) (*model.AuditResult, error) {
	started := time.Now()
	if opts.Niche == "" {
		opts.Niche = "local-service"
	}
	if !supportedNiches[opts.Niche] {
		return nil, fmt.Errorf("unsupported niche: %s", opts.Niche)
	}
	fetched, err := fetch.Fetch(ctx, rawURL, fetch.Options{TimeoutSeconds: opts.TimeoutSeconds, MaxBytes: opts.MaxBytes, UserAgent: opts.UserAgent})
	if err != nil {
		return nil, err
	}
	extracted, err := extract.FromHTML(string(fetched.Body), opts.Niche)
	if err != nil {
		return nil, err
	}
	points, grade := score.Calculate(extracted.Signals)
	result := &model.AuditResult{
		URL:             fetched.URL,
		Niche:           opts.Niche,
		Score:           points,
		Grade:           grade,
		Title:           extracted.Title,
		MetaDescription: extracted.MetaDescription,
		H1:              extracted.H1,
		H2s:             extracted.H2s,
		Signals:         extracted.Signals,
		Issues:          score.Issues(extracted.Signals),
		Recommendations: score.Recommendations(extracted.Signals),
		OutreachAngle:   score.OutreachAngle(points),
		Diagnostics: model.Diagnostics{
			StatusCode:       fetched.StatusCode,
			ContentType:      fetched.ContentType,
			BytesRead:        len(fetched.Body),
			MatchedTerms:     extracted.MatchedTerms,
			ProcessingMillis: time.Since(started).Milliseconds(),
		},
	}
	return result, nil
}

// RunHTML audits a fixed HTML string. It is used by tests and future wrappers.
func RunHTML(rawURL, html string, opts Options) (*model.AuditResult, error) {
	started := time.Now()
	if opts.Niche == "" {
		opts.Niche = "local-service"
	}
	normalized, err := fetch.URLString(rawURL)
	if err != nil {
		return nil, err
	}
	if !supportedNiches[opts.Niche] {
		return nil, fmt.Errorf("unsupported niche: %s", opts.Niche)
	}
	extracted, err := extract.FromHTML(html, opts.Niche)
	if err != nil {
		return nil, err
	}
	points, grade := score.Calculate(extracted.Signals)
	return &model.AuditResult{
		URL:             normalized,
		Niche:           opts.Niche,
		Score:           points,
		Grade:           grade,
		Title:           extracted.Title,
		MetaDescription: extracted.MetaDescription,
		H1:              extracted.H1,
		H2s:             extracted.H2s,
		Signals:         extracted.Signals,
		Issues:          score.Issues(extracted.Signals),
		Recommendations: score.Recommendations(extracted.Signals),
		OutreachAngle:   score.OutreachAngle(points),
		Diagnostics: model.Diagnostics{
			ContentType:      "text/html; fixture",
			BytesRead:        len([]byte(html)),
			MatchedTerms:     extracted.MatchedTerms,
			ProcessingMillis: time.Since(started).Milliseconds(),
		},
	}, nil
}

// SupportedNiches returns a stable display list.
func SupportedNiches() string {
	return strings.Join([]string{"legal", "local-service"}, ", ")
}

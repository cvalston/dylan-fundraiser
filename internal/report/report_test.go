package report

import (
	"encoding/json"
	"strings"
	"testing"

	engine "siteintel/internal/audit"
)

const fixture = `<!doctype html><html><head><title>West Orange Personal Injury Lawyers</title><meta name="description" content="Serving New Jersey injury victims."></head><body><h1>Personal Injury Law Firm Serving West Orange</h1><a href="/contact">Contact Us</a><p>Call (973) 555-1212 for a free consultation. Years of experience. Testimonials. Car accident help.</p><form></form></body></html>`

func TestMarkdownRendering(t *testing.T) {
	result, err := engine.RunHTML("https://example.com", fixture, engine.Options{Niche: "legal"})
	if err != nil {
		t.Fatalf("audit failed: %v", err)
	}
	out := Markdown(*result)
	if !strings.Contains(out, "# SiteIntel Audit") || !strings.Contains(out, "## Score") {
		t.Fatalf("markdown missing expected sections: %s", out)
	}
}

func TestJSONRendering(t *testing.T) {
	result, err := engine.RunHTML("https://example.com", fixture, engine.Options{Niche: "legal"})
	if err != nil {
		t.Fatalf("audit failed: %v", err)
	}
	out, err := JSON(*result)
	if err != nil {
		t.Fatalf("json failed: %v", err)
	}
	var decoded map[string]any
	if err := json.Unmarshal([]byte(out), &decoded); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if decoded["url"] != "https://example.com" {
		t.Fatalf("unexpected url: %v", decoded["url"])
	}
}

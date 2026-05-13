package report

import (
	"fmt"
	"strings"

	"siteintel/internal/model"
	"siteintel/internal/score"
)

// Markdown renders compact deterministic markdown.
func Markdown(result model.AuditResult) string {
	var b strings.Builder
	b.WriteString("# SiteIntel Audit\n\n")
	b.WriteString("## Summary\n")
	b.WriteString(summary(result.Score) + "\n\n")
	b.WriteString("## Score\n")
	b.WriteString(fmt.Sprintf("%d/100 — %s\n\n", result.Score, result.Grade))
	b.WriteString("## Detected Strengths\n")
	writeBullets(&b, score.Strengths(result.Signals))
	b.WriteString("\n## Priority Issues\n")
	writeBullets(&b, result.Issues)
	b.WriteString("\n## Recommendations\n")
	writeBullets(&b, result.Recommendations)
	b.WriteString("\n## Outreach Angle\n")
	b.WriteString(result.OutreachAngle + "\n\n")
	b.WriteString("## Extracted Signals\n")
	b.WriteString(fmt.Sprintf("- Title: %s\n", fallback(result.Title)))
	b.WriteString(fmt.Sprintf("- Meta description: %s\n", fallback(result.MetaDescription)))
	b.WriteString(fmt.Sprintf("- H1: %s\n", fallback(result.H1)))
	if len(result.H2s) > 0 {
		b.WriteString(fmt.Sprintf("- H2s: %s\n", strings.Join(result.H2s, " | ")))
	} else {
		b.WriteString("- H2s: Not detected\n")
	}
	b.WriteString(fmt.Sprintf("- Phone: %t\n", result.Signals.HasPhone))
	b.WriteString(fmt.Sprintf("- Contact link: %t\n", result.Signals.HasContactLink))
	b.WriteString(fmt.Sprintf("- Consultation CTA: %t\n", result.Signals.HasConsultationCTA))
	b.WriteString(fmt.Sprintf("- Form: %t\n", result.Signals.HasForm))
	b.WriteString(fmt.Sprintf("- Trust indicators: %t\n", result.Signals.HasTrustIndicators))
	b.WriteString(fmt.Sprintf("- Practice areas/services: %t\n", result.Signals.HasPracticeAreas))
	b.WriteString(fmt.Sprintf("- Location signal: %t\n", result.Signals.HasLocationSignal))
	b.WriteString(fmt.Sprintf("- Schema: %t\n", result.Signals.HasSchema))
	b.WriteString("\n## Diagnostics\n")
	b.WriteString(fmt.Sprintf("- URL: %s\n", result.URL))
	b.WriteString(fmt.Sprintf("- Niche: %s\n", result.Niche))
	if result.Diagnostics.StatusCode != 0 {
		b.WriteString(fmt.Sprintf("- HTTP status: %d\n", result.Diagnostics.StatusCode))
	}
	if result.Diagnostics.ContentType != "" {
		b.WriteString(fmt.Sprintf("- Content type: %s\n", result.Diagnostics.ContentType))
	}
	if result.Diagnostics.BytesRead != 0 {
		b.WriteString(fmt.Sprintf("- Bytes read: %d\n", result.Diagnostics.BytesRead))
	}
	if result.Diagnostics.ProcessingMillis != 0 {
		b.WriteString(fmt.Sprintf("- Processing ms: %d\n", result.Diagnostics.ProcessingMillis))
	}
	if len(result.Diagnostics.MatchedTerms) > 0 {
		b.WriteString(fmt.Sprintf("- Matched terms: %s\n", strings.Join(result.Diagnostics.MatchedTerms, ", ")))
	}
	return b.String()
}

func writeBullets(b *strings.Builder, items []string) {
	if len(items) == 0 {
		b.WriteString("- None\n")
		return
	}
	for _, item := range items {
		b.WriteString("- " + item + "\n")
	}
}

func fallback(value string) string {
	if strings.TrimSpace(value) == "" {
		return "Not detected"
	}
	return value
}

func summary(score int) string {
	switch {
	case score >= 80:
		return "Website appears conversion-ready with opportunities for incremental gains."
	case score >= 70:
		return "Website appears credible but has several conversion weaknesses."
	case score >= 60:
		return "Website has meaningful conversion gaps that may reduce lead generation."
	default:
		return "Website shows significant audit gaps and likely missed lead opportunities."
	}
}

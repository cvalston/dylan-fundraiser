package extract

import "testing"

const fixtureHTML = `<!doctype html>
<html><head>
<title>West Orange Personal Injury Lawyers</title>
<meta name="description" content="Serving New Jersey injury victims with award winning attorneys.">
<script type="application/ld+json">{"@context":"https://schema.org","@type":"LegalService"}</script>
</head><body>
<nav><a href="/contact">Contact Us</a><a href="/about">About Our Team</a><a href="https://facebook.com/example">Facebook</a><a href="/privacy-policy">Privacy Policy</a></nav>
<h1>Personal Injury Law Firm Serving West Orange</h1>
<h2>Car Accident Lawyers</h2>
<p>Call today for a free consultation. Our attorneys have 25 years of experience and million dollar results.</p>
<p>Client reviews and testimonials show trusted support across Essex County.</p>
<p>Call (973) 555-1212.</p>
<form><input name="name"><textarea name="message"></textarea></form>
</body></html>`

func TestSignalExtraction(t *testing.T) {
	result, err := FromHTML(fixtureHTML, "legal")
	if err != nil {
		t.Fatalf("extract failed: %v", err)
	}
	if result.Title != "West Orange Personal Injury Lawyers" {
		t.Fatalf("unexpected title: %q", result.Title)
	}
	checks := map[string]bool{
		"phone":       result.Signals.HasPhone,
		"cta":         result.Signals.HasConsultationCTA,
		"trust":       result.Signals.HasTrustIndicators,
		"practice":    result.Signals.HasPracticeAreas,
		"location":    result.Signals.HasLocationSignal,
		"schema":      result.Signals.HasSchema,
		"social":      result.Signals.HasSocialLinks,
		"privacy":     result.Signals.HasPrivacyPolicy,
		"headline":    result.Signals.HasStrongHeadline,
		"primary_cta": result.Signals.HasClearPrimaryCTA,
	}
	for name, ok := range checks {
		if !ok {
			t.Fatalf("expected %s signal", name)
		}
	}
}

func TestPhoneDetection(t *testing.T) {
	if !HasPhone("Call 212-555-0199 today") {
		t.Fatal("expected phone detection")
	}
	if HasPhone("No number here") {
		t.Fatal("did not expect phone detection")
	}
}

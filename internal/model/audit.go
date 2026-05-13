package model

// Signals captures deterministic homepage audit signals.
type Signals struct {
	HasPhone           bool `json:"has_phone"`
	HasContactLink     bool `json:"has_contact_link"`
	HasConsultationCTA bool `json:"has_consultation_cta"`
	HasForm            bool `json:"has_form"`
	HasTestimonials    bool `json:"has_testimonials"`
	HasReviews         bool `json:"has_reviews"`
	HasTrustIndicators bool `json:"has_trust_indicators"`
	HasAboutOrTeam     bool `json:"has_about_or_team"`
	HasPracticeAreas   bool `json:"has_practice_areas"`
	HasLocationSignal  bool `json:"has_location_signal"`
	HasSchema          bool `json:"has_schema"`
	HasSocialLinks     bool `json:"has_social_links"`
	HasPrivacyPolicy   bool `json:"has_privacy_policy"`
	HasStrongHeadline  bool `json:"has_strong_headline"`
	HasClearPrimaryCTA bool `json:"has_clear_primary_cta"`
}

// Diagnostics contains compact, non-sensitive processing details.
type Diagnostics struct {
	StatusCode       int      `json:"status_code,omitempty"`
	ContentType      string   `json:"content_type,omitempty"`
	BytesRead        int      `json:"bytes_read,omitempty"`
	MatchedTerms     []string `json:"matched_terms,omitempty"`
	Warnings         []string `json:"warnings,omitempty"`
	ProcessingMillis int64    `json:"processing_millis,omitempty"`
}

// AuditResult is the top-level deterministic audit output.
type AuditResult struct {
	URL             string      `json:"url"`
	Niche           string      `json:"niche"`
	Score           int         `json:"score"`
	Grade           string      `json:"grade"`
	Title           string      `json:"title,omitempty"`
	MetaDescription string      `json:"meta_description,omitempty"`
	H1              string      `json:"h1,omitempty"`
	H2s             []string    `json:"h2s,omitempty"`
	Signals         Signals     `json:"signals"`
	Issues          []string    `json:"issues"`
	Recommendations []string    `json:"recommendations"`
	OutreachAngle   string      `json:"outreach_angle"`
	Diagnostics     Diagnostics `json:"diagnostics"`
}

package score

import "siteintel/internal/model"

// Weighted signal values, intentionally explicit and auditable.
var weights = map[string]int{
	"phone":                8,
	"contact_link":         8,
	"consultation_cta":     10,
	"form":                 8,
	"reviews_testimonials": 10,
	"trust_indicators":     10,
	"about_team":           7,
	"practice_areas":       8,
	"location_signal":      8,
	"schema":               6,
	"social_links":         4,
	"privacy_policy":       3,
	"strong_headline":      5,
	"clear_primary_cta":    5,
}

// Calculate returns a capped score and deterministic grade.
func Calculate(signals model.Signals) (int, string) {
	score := 0
	if signals.HasPhone {
		score += weights["phone"]
	}
	if signals.HasContactLink {
		score += weights["contact_link"]
	}
	if signals.HasConsultationCTA {
		score += weights["consultation_cta"]
	}
	if signals.HasForm {
		score += weights["form"]
	}
	if signals.HasTestimonials || signals.HasReviews {
		score += weights["reviews_testimonials"]
	}
	if signals.HasTrustIndicators {
		score += weights["trust_indicators"]
	}
	if signals.HasAboutOrTeam {
		score += weights["about_team"]
	}
	if signals.HasPracticeAreas {
		score += weights["practice_areas"]
	}
	if signals.HasLocationSignal {
		score += weights["location_signal"]
	}
	if signals.HasSchema {
		score += weights["schema"]
	}
	if signals.HasSocialLinks {
		score += weights["social_links"]
	}
	if signals.HasPrivacyPolicy {
		score += weights["privacy_policy"]
	}
	if signals.HasStrongHeadline {
		score += weights["strong_headline"]
	}
	if signals.HasClearPrimaryCTA {
		score += weights["clear_primary_cta"]
	}
	if score > 100 {
		score = 100
	}
	return score, Grade(score)
}

// Grade maps scores to deterministic grade labels.
func Grade(score int) string {
	switch {
	case score >= 90:
		return "Excellent"
	case score >= 80:
		return "Strong"
	case score >= 70:
		return "Good but improvable"
	case score >= 60:
		return "Needs improvement"
	default:
		return "High opportunity"
	}
}

// Issues creates compact issue statements from weak or missing signals.
func Issues(signals model.Signals) []string {
	var issues []string
	if !signals.HasPhone {
		issues = append(issues, "No clear phone number detected.")
	}
	if !signals.HasContactLink {
		issues = append(issues, "No clear contact page link detected.")
	}
	if !signals.HasConsultationCTA {
		issues = append(issues, "No strong consultation CTA detected.")
	}
	if !signals.HasForm {
		issues = append(issues, "No contact form detected.")
	}
	if !signals.HasTestimonials && !signals.HasReviews {
		issues = append(issues, "No testimonials or reviews detected.")
	}
	if !signals.HasTrustIndicators {
		issues = append(issues, "Weak trust indicators detected.")
	}
	if !signals.HasPracticeAreas {
		issues = append(issues, "No clear practice area structure detected.")
	}
	if !signals.HasLocationSignal {
		issues = append(issues, "No local presence signal detected.")
	}
	if !signals.HasSchema {
		issues = append(issues, "No structured data detected.")
	}
	if !signals.HasStrongHeadline {
		issues = append(issues, "Headline does not clearly communicate the business focus.")
	}
	if !signals.HasClearPrimaryCTA {
		issues = append(issues, "Weak primary CTA detected.")
	}
	return issues
}

// Recommendations creates deterministic improvement templates.
func Recommendations(signals model.Signals) []string {
	var recs []string
	if !signals.HasConsultationCTA || !signals.HasClearPrimaryCTA {
		recs = append(recs, "Add a clear above-the-fold consultation CTA such as “Schedule a Free Consultation” or “Request a Case Review.”")
	}
	if !signals.HasTrustIndicators {
		recs = append(recs, "Add visible credibility signals such as years of experience, awards, notable results, testimonials, or client review highlights.")
	}
	if !signals.HasTestimonials && !signals.HasReviews {
		recs = append(recs, "Add testimonials and review highlights near primary conversion sections.")
	}
	if !signals.HasPhone {
		recs = append(recs, "Add a highly visible phone number in the header and mobile navigation.")
	}
	if !signals.HasLocationSignal {
		recs = append(recs, "Add stronger geographic relevance throughout the homepage and metadata.")
	}
	if !signals.HasPracticeAreas {
		recs = append(recs, "Add a scannable practice area or service section with links to key offerings.")
	}
	if !signals.HasSchema {
		recs = append(recs, "Add LocalBusiness or LegalService structured data to improve machine readability.")
	}
	return recs
}

// OutreachAngle returns a score-band-specific outreach angle.
func OutreachAngle(score int) string {
	switch {
	case score >= 90:
		return "Focus on optimization and competitive positioning."
	case score >= 80:
		return "Focus on incremental conversion gains."
	case score >= 70:
		return "Focus on visible conversion gaps."
	case score >= 60:
		return "Focus on meaningful missed opportunity."
	default:
		return "Focus on significant underperformance and missed lead generation."
	}
}

// Strengths summarizes positive signals for compact reports.
func Strengths(signals model.Signals) []string {
	var strengths []string
	if signals.HasPracticeAreas {
		strengths = append(strengths, "Practice areas detected")
	}
	if signals.HasForm {
		strengths = append(strengths, "Contact form detected")
	}
	if signals.HasPhone {
		strengths = append(strengths, "Phone number detected")
	}
	if signals.HasConsultationCTA {
		strengths = append(strengths, "Consultation CTA detected")
	}
	if signals.HasTrustIndicators {
		strengths = append(strengths, "Trust indicators detected")
	}
	if signals.HasTestimonials || signals.HasReviews {
		strengths = append(strengths, "Reviews or testimonials detected")
	}
	if signals.HasLocationSignal {
		strengths = append(strengths, "Local presence signals detected")
	}
	if len(strengths) == 0 {
		strengths = append(strengths, "No major conversion strengths detected")
	}
	return strengths
}

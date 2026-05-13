package score

import (
	"testing"

	"siteintel/internal/model"
)

func TestCalculateDeterministicScore(t *testing.T) {
	signals := model.Signals{HasPhone: true, HasContactLink: true, HasConsultationCTA: true, HasForm: true, HasTestimonials: true, HasTrustIndicators: true, HasAboutOrTeam: true, HasPracticeAreas: true, HasLocationSignal: true, HasSchema: true, HasSocialLinks: true, HasPrivacyPolicy: true, HasStrongHeadline: true, HasClearPrimaryCTA: true}
	score, grade := Calculate(signals)
	if score != 100 {
		t.Fatalf("expected capped score 100, got %d", score)
	}
	if grade != "Excellent" {
		t.Fatalf("expected Excellent, got %q", grade)
	}
}

func TestIssuesAndRecommendations(t *testing.T) {
	issues := Issues(model.Signals{})
	recs := Recommendations(model.Signals{})
	if len(issues) == 0 || len(recs) == 0 {
		t.Fatal("expected issues and recommendations for empty signals")
	}
}

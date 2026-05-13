package extract

import (
	"regexp"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"siteintel/internal/model"
)

var phonePattern = regexp.MustCompile(`(?i)(?:\+?1[\s.-]?)?(?:\([2-9][0-9]{2}\)|[2-9][0-9]{2})[\s.-]?[0-9]{3}[\s.-]?[0-9]{4}`)

var ctaTerms = []string{"free consultation", "schedule a consultation", "request a consultation", "case evaluation", "call today", "contact us", "speak with an attorney", "get help now"}
var trustTerms = []string{"years of experience", "award", "recognized", "top rated", "super lawyers", "million", "settlement", "verdict", "results", "client reviews", "testimonials"}
var practiceAreaTerms = []string{"personal injury", "car accident", "workers compensation", "medical malpractice", "criminal defense", "family law", "employment law", "real estate", "estate planning"}
var locationTerms = []string{"serving", "located in", " near ", "new jersey", "newark", "west orange", "essex county", "nyc", "new york"}
var socialDomains = []string{"facebook.com", "instagram.com", "linkedin.com", "x.com", "twitter.com", "youtube.com", "tiktok.com"}

// Result is deterministic extraction output before scoring.
type Result struct {
	Title           string
	MetaDescription string
	H1              string
	H2s             []string
	Signals         model.Signals
	MatchedTerms    []string
}

// FromHTML extracts visible website signals from homepage HTML.
func FromHTML(html string, niche string) (*Result, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}
	doc.Find("script,style,noscript,svg").Remove()

	title := clean(doc.Find("title").First().Text())
	metaDesc, _ := doc.Find(`meta[name="description"], meta[property="og:description"]`).First().Attr("content")
	metaDesc = clean(metaDesc)
	h1 := clean(doc.Find("h1").First().Text())
	var h2s []string
	doc.Find("h2").Each(func(_ int, s *goquery.Selection) {
		text := clean(s.Text())
		if text != "" && len(h2s) < 10 {
			h2s = append(h2s, text)
		}
	})

	bodyText := " " + strings.ToLower(clean(doc.Find("body").Text())) + " "
	allText := " " + strings.ToLower(clean(title+" "+metaDesc+" "+h1+" "+strings.Join(h2s, " ")+" "+bodyText)) + " "

	linksText := strings.ToLower(linkTextAndHrefs(doc))
	buttonText := strings.ToLower(clean(doc.Find("a,button,input[type=submit]").Text()))

	signals := model.Signals{}
	signals.HasPhone = phonePattern.MatchString(allText)
	signals.HasContactLink = strings.Contains(linksText, "contact") || strings.Contains(linksText, "/contact")
	signals.HasConsultationCTA = containsAny(allText, ctaTerms)
	signals.HasForm = doc.Find("form, input, textarea, select").Length() > 0
	signals.HasTestimonials = strings.Contains(allText, "testimonial") || strings.Contains(allText, "what our clients say")
	signals.HasReviews = strings.Contains(allText, "review") || strings.Contains(allText, "google rating") || strings.Contains(allText, "stars")
	signals.HasTrustIndicators = containsAny(allText, trustTerms)
	signals.HasAboutOrTeam = strings.Contains(linksText, "about") || strings.Contains(linksText, "team") || strings.Contains(allText, "our team") || strings.Contains(allText, "attorneys")
	signals.HasPracticeAreas = containsAny(allText, practiceAreaTerms) || strings.Contains(allText, "practice areas") || strings.Contains(allText, "services")
	signals.HasLocationSignal = containsAny(allText, locationTerms)
	signals.HasSchema = doc.Find(`script[type="application/ld+json"], [itemscope], [itemtype*="schema.org"]`).Length() > 0 || strings.Contains(html, "schema.org")
	signals.HasSocialLinks = containsAny(linksText, socialDomains)
	signals.HasPrivacyPolicy = strings.Contains(linksText, "privacy") || strings.Contains(allText, "privacy policy")
	signals.HasStrongHeadline = strongHeadline(h1, title, niche)
	signals.HasClearPrimaryCTA = primaryCTA(buttonText + " " + linksText)

	matched := matchedTerms(allText, append(append(append([]string{}, ctaTerms...), trustTerms...), append(practiceAreaTerms, locationTerms...)...))
	return &Result{Title: title, MetaDescription: metaDesc, H1: h1, H2s: h2s, Signals: signals, MatchedTerms: matched}, nil
}

func clean(value string) string {
	return strings.Join(strings.Fields(value), " ")
}

func containsAny(haystack string, needles []string) bool {
	for _, needle := range needles {
		if strings.Contains(haystack, needle) {
			return true
		}
	}
	return false
}

func linkTextAndHrefs(doc *goquery.Document) string {
	var parts []string
	doc.Find("a").Each(func(_ int, s *goquery.Selection) {
		parts = append(parts, s.Text())
		if href, ok := s.Attr("href"); ok {
			parts = append(parts, href)
		}
	})
	return clean(strings.Join(parts, " "))
}

func strongHeadline(h1, title, niche string) bool {
	text := strings.ToLower(h1 + " " + title)
	if len(strings.Fields(h1)) < 4 {
		return false
	}
	if niche == "legal" {
		return strings.Contains(text, "law") || strings.Contains(text, "attorney") || strings.Contains(text, "lawyer") || containsAny(text, practiceAreaTerms)
	}
	return strings.Contains(text, "service") || strings.Contains(text, "help") || strings.Contains(text, "local") || len(strings.Fields(h1)) >= 6
}

func primaryCTA(text string) bool {
	return containsAny(text, ctaTerms) || strings.Contains(text, "schedule") || strings.Contains(text, "book") || strings.Contains(text, "call")
}

func matchedTerms(text string, terms []string) []string {
	seen := map[string]bool{}
	for _, term := range terms {
		if strings.Contains(text, term) {
			seen[strings.TrimSpace(term)] = true
		}
	}
	out := make([]string, 0, len(seen))
	for term := range seen {
		if term != "" {
			out = append(out, term)
		}
	}
	sort.Strings(out)
	return out
}

// HasPhone exposes phone detection for tests and future packages.
func HasPhone(text string) bool { return phonePattern.MatchString(text) }

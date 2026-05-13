# siteintel

`siteintel` is a deterministic Go CLI for auditing local business and law firm homepages. It fetches public HTML, removes noisy content, extracts conversion and trust signals with rules, scores the site with explicit weights, and emits compact reports for humans, AI agents, future MCP tools, and API wrappers.

## Architecture philosophy

The v1 processing path is intentionally deterministic:

```text
Website URL -> deterministic CLI -> compact structured report -> optional AI reasoning
```

`siteintel` does **not** use an LLM for fetching, parsing, extraction, classification, scoring, issue generation, recommendation generation, or report rendering. The same input HTML produces the same score and report fields.

The code separates concerns for maintainability and future orchestration:

- `cmd/`: Cobra command wiring only.
- `internal/fetch/`: URL validation, HTTP fetching, timeouts, body limits, content-type checks.
- `internal/extract/`: HTML parsing, script/style removal, deterministic signal extraction.
- `internal/score/`: explicit weighted scoring, grades, issues, recommendations, outreach angle.
- `internal/report/`: Markdown and JSON rendering.
- `internal/model/`: stable output schema.
- `internal/audit/`: orchestration layer usable by CLI, tests, MCP wrappers, or APIs.

## Install

```bash
go mod download
go install .
```

Or run without installing:

```bash
go run . audit https://example.com --niche legal
```

## Usage

```bash
siteintel audit <url> [flags]
```

Common examples:

```bash
siteintel audit https://examplelawfirm.com --niche legal
siteintel audit https://examplelawfirm.com --niche legal --format json
siteintel audit https://example.com --timeout 10 --max-bytes 1000000
```

Supported flags:

- `--niche string`: `legal` or `local-service` (default `local-service`)
- `--format string`: `markdown` or `json` (default `markdown`)
- `--timeout int`: HTTP timeout in seconds (default `15`)
- `--max-bytes int`: maximum response size to process (default `1500000`)
- `--user-agent string`: HTTP user agent (default `siteintel/0.1 (+https://example.com)`)
- `--verbose`: show additional diagnostics

## Output examples

Markdown output is compact and decision-ready:

```markdown
# SiteIntel Audit

## Summary
Website appears credible but has several conversion weaknesses.

## Score
71/100 — Good but improvable

## Detected Strengths
- Practice areas detected
- Contact form detected
- Phone number detected

## Priority Issues
- Weak primary CTA detected.
- No testimonials or reviews detected.

## Recommendations
- Add a clear above-the-fold consultation CTA such as “Schedule a Free Consultation” or “Request a Case Review.”
- Add testimonials and review highlights near primary conversion sections.

## Outreach Angle
Focus on visible conversion gaps.
```

JSON output contains the same important fields in a machine-readable schema suitable for future MCP and API use:

```bash
siteintel audit https://examplelawfirm.com --niche legal --format json
```

## Deterministic scoring

Scores are built from explicit signal weights and capped at 100:

- Phone number: 8
- Contact link: 8
- Consultation CTA: 10
- Form: 8
- Reviews/testimonials: 10
- Trust indicators: 10
- About/team: 7
- Practice areas/services: 8
- Location signal: 8
- Schema: 6
- Social links: 4
- Privacy policy: 3
- Strong headline: 5
- Clear primary CTA: 5

Grade ranges:

- `90-100`: Excellent
- `80-89`: Strong
- `70-79`: Good but improvable
- `60-69`: Needs improvement
- `<60`: High opportunity

## Security and scope

V1 analyzes only publicly visible homepage HTML. It does not execute remote JavaScript, run browser automation, bypass authentication, bypass paywalls, or collect hidden/private data.

## MCP extensibility notes

The `internal/audit` package exposes orchestration functions that return a stable `AuditResult` model. A future MCP server can wrap that package directly and expose deterministic `audit_url` or `audit_html` tools without coupling to Cobra command code. The JSON renderer already produces compact machine-readable output for agent workflows.

## Future roadmap

Planned extensions include:

- MCP wrapper
- REST API
- multi-page crawling
- DataForSEO enrichment
- Google Business Profile enrichment
- screenshot analysis
- Remotion brief generation
- outreach email generation
- Supabase persistence
- dashboards
- batch analysis
- optional AI reasoning layer over compact deterministic output
- workflow orchestration

# siteintel

Use `siteintel` when an agent workflow needs a compact, deterministic homepage audit for a local business or law firm website.

## Purpose

`siteintel` preprocesses public homepage HTML into a structured audit report with deterministic extraction, weighted scoring, priority issues, recommendations, and an outreach angle.

## Deterministic constraints

Do not use LLMs for v1 extraction, classification, scoring, recommendations, or report rendering. Use only rule-based parsing, keyword dictionaries, regular expressions, weighted scoring, and fixed templates.

## CLI examples

```bash
go run . audit https://example.com --niche legal
```

```bash
go run . audit https://example.com --niche legal --format json
```

## Agent integration notes

Prefer JSON output for MCP, API, or workflow orchestration. Prefer Markdown output for quick human review. The `internal/audit` package is the intended integration seam for future wrappers.

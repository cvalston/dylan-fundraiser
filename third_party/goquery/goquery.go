package goquery

import (
	"io"
	"regexp"
	"strings"
)

type Document struct{ html string }
type Selection struct {
	nodes []node
	doc   *Document
}
type node struct {
	tag   string
	attrs map[string]string
	inner string
	full  string
}

func NewDocumentFromReader(r io.Reader) (*Document, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return &Document{html: string(b)}, nil
}

func (d *Document) Find(selector string) *Selection { return find(d.html, selector, d) }
func (s *Selection) Find(selector string) *Selection {
	var all []node
	for _, n := range s.nodes {
		all = append(all, find(n.inner, selector, s.doc).nodes...)
	}
	return &Selection{nodes: all, doc: s.doc}
}
func (s *Selection) Remove() {
	if s.doc == nil {
		return
	}
	for _, n := range s.nodes {
		s.doc.html = strings.ReplaceAll(s.doc.html, n.full, "")
	}
}
func (s *Selection) First() *Selection {
	if len(s.nodes) == 0 {
		return &Selection{doc: s.doc}
	}
	return &Selection{nodes: s.nodes[:1], doc: s.doc}
}
func (s *Selection) Text() string {
	parts := make([]string, 0, len(s.nodes))
	for _, n := range s.nodes {
		if n.tag == "input" {
			if v := n.attrs["value"]; v != "" {
				parts = append(parts, v)
			}
			continue
		}
		parts = append(parts, stripTags(n.inner))
	}
	return strings.Join(parts, " ")
}
func (s *Selection) Attr(name string) (string, bool) {
	if len(s.nodes) == 0 {
		return "", false
	}
	v, ok := s.nodes[0].attrs[strings.ToLower(name)]
	return v, ok
}
func (s *Selection) Each(fn func(int, *Selection)) {
	for i := range s.nodes {
		fn(i, &Selection{nodes: s.nodes[i : i+1], doc: s.doc})
	}
}
func (s *Selection) Length() int { return len(s.nodes) }

func find(html, selector string, doc *Document) *Selection {
	var out []node
	for _, sel := range strings.Split(selector, ",") {
		sel = strings.TrimSpace(sel)
		switch {
		case sel == "title" || sel == "h1" || sel == "h2" || sel == "body" || sel == "form" || sel == "textarea" || sel == "select" || sel == "button" || sel == "a":
			out = append(out, tagNodes(html, sel)...)
		case sel == "input" || sel == "input[type=submit]":
			for _, n := range voidNodes(html, "input") {
				if sel == "input" || strings.EqualFold(n.attrs["type"], "submit") {
					out = append(out, n)
				}
			}
		case strings.HasPrefix(sel, "meta"):
			for _, n := range voidNodes(html, "meta") {
				name := strings.ToLower(n.attrs["name"])
				prop := strings.ToLower(n.attrs["property"])
				if name == "description" || prop == "og:description" {
					out = append(out, n)
				}
			}
		case strings.HasPrefix(sel, "script[type="):
			for _, n := range tagNodes(html, "script") {
				if strings.EqualFold(n.attrs["type"], "application/ld+json") {
					out = append(out, n)
				}
			}
		case sel == "[itemscope]":
			if strings.Contains(strings.ToLower(html), "itemscope") {
				out = append(out, node{tag: "itemscope", full: "itemscope"})
			}
		case strings.HasPrefix(sel, "[itemtype"):
			if strings.Contains(strings.ToLower(html), "schema.org") {
				out = append(out, node{tag: "itemtype", full: "schema.org"})
			}
		}
	}
	return &Selection{nodes: out, doc: doc}
}

func tagNodes(html, tag string) []node {
	re := regexp.MustCompile(`(?is)<` + regexp.QuoteMeta(tag) + `\b([^>]*)>(.*?)</` + regexp.QuoteMeta(tag) + `>`)
	matches := re.FindAllStringSubmatch(html, -1)
	out := make([]node, 0, len(matches))
	for _, m := range matches {
		out = append(out, node{tag: tag, attrs: attrs(m[1]), inner: m[2], full: m[0]})
	}
	return out
}

func voidNodes(html, tag string) []node {
	re := regexp.MustCompile(`(?is)<` + regexp.QuoteMeta(tag) + `\b([^>]*)/?>`)
	matches := re.FindAllStringSubmatch(html, -1)
	out := make([]node, 0, len(matches))
	for _, m := range matches {
		out = append(out, node{tag: tag, attrs: attrs(m[1]), full: m[0]})
	}
	return out
}

func attrs(raw string) map[string]string {
	out := map[string]string{}
	re := regexp.MustCompile(`(?is)([a-zA-Z_:][-a-zA-Z0-9_:.]*)\s*=\s*("([^"]*)"|'([^']*)'|([^\s>]+))`)
	for _, m := range re.FindAllStringSubmatch(raw, -1) {
		v := m[3]
		if v == "" {
			v = m[4]
		}
		if v == "" {
			v = m[5]
		}
		out[strings.ToLower(m[1])] = v
	}
	return out
}

func stripTags(s string) string { return regexp.MustCompile(`(?is)<[^>]+>`).ReplaceAllString(s, " ") }

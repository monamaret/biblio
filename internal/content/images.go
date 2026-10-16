// Copyright 2026 Mona Maret
// SPDX-License-Identifier: Apache-2.0

package content

import (
	"path"
	"sort"
	"strings"

	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

// RewriteImages points every inline image whose destination resolves to a
// converted plate at that plate's URL. dir is the content-relative folder
// of the file the body belongs to (relative destinations resolve against
// it); plates maps content-relative image paths to plate URLs. Only
// destinations goldmark parses as images are touched — never URLs it
// doesn't know, never text inside code (specs/004-quarterly-site-type
// FR-028; research.md R5).
func RewriteImages(body, dir string, plates map[string]string) string {
	if len(plates) == 0 {
		return body
	}
	src := []byte(body)
	doc := notesMarkdown.Parser().Parse(text.NewReader(src))

	targets := map[string]string{}
	_ = gast.Walk(doc, func(n gast.Node, entering bool) (gast.WalkStatus, error) {
		if img, ok := n.(*gast.Image); ok && entering {
			dest := string(img.Destination)
			if url, ok := plates[path.Clean(path.Join(dir, dest))]; ok && !strings.Contains(dest, "://") {
				targets[dest] = url
			}
		}
		return gast.WalkContinue, nil
	})
	if len(targets) == 0 {
		return body
	}

	code := codeRanges(doc, src)
	type hit struct {
		start, end int
		url        string
	}
	var hits []hit
	for dest, url := range targets {
		needle := "](" + dest
		for i := 0; ; {
			j := strings.Index(body[i:], needle)
			if j < 0 {
				break
			}
			start := i + j + 2
			end := start + len(dest)
			if !inRanges(code, start) && end <= len(body) && (end == len(body) || body[end] == ')' || body[end] == ' ') {
				hits = append(hits, hit{start, end, url})
			}
			i = end
		}
	}
	sort.Slice(hits, func(a, b int) bool { return hits[a].start < hits[b].start })

	var b strings.Builder
	last := 0
	for _, h := range hits {
		if h.start < last {
			continue
		}
		b.WriteString(body[last:h.start])
		b.WriteString(h.url)
		last = h.end
	}
	b.WriteString(body[last:])
	return b.String()
}

package content

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/yuin/goldmark"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/text"
)

// Note is one footnote or margin note of a quarterly article, as stored in
// content-index.json (specs/004-quarterly-site-type data-model.md "Note").
type Note struct {
	// ID is "fn-<label>" for a footnote, "m1", "m2", … for margin notes.
	ID string `json:"id"`
	// Label is the author's footnote label; empty for margin notes.
	Label string `json:"label,omitempty"`
	// Number is the footnote's display number (1-based, in reference
	// order); zero for margin notes.
	Number int `json:"number,omitempty"`
	// Text is the note's inline Markdown.
	Text string `json:"text"`
	// RefOffset is the Unicode code point offset of the note's reference
	// marker ("[^label]" or "{>>") in the article body.
	RefOffset int `json:"ref_offset"`
}

var (
	notesMarkdown = goldmark.New(goldmark.WithExtensions(extension.Footnote))
	footnoteRef   = regexp.MustCompile(`\[\^([^\]\s]+)\]`)
)

// ExtractNotes returns an article body's footnotes (in order of first
// reference) and margin notes (in document order), plus warnings for
// problems that never fail the build: a reference with no definition, a
// definition never referenced, an unclosed {>>. Markers inside code spans
// and code blocks are never notes (FR-016; research.md R2).
//
// Footnote definitions come from goldmark's footnote extension; reference
// positions and margin notes come from scanning the source outside the
// code ranges goldmark reports, because goldmark keeps no source position
// for a footnote reference.
func ExtractNotes(body string) (footnotes, marginNotes []Note, warnings []string) {
	src := []byte(body)
	doc := notesMarkdown.Parser().Parse(text.NewReader(src))

	code := codeRanges(doc, src)
	defs := map[string]string{}
	var defOrder []string
	_ = gast.Walk(doc, func(n gast.Node, entering bool) (gast.WalkStatus, error) {
		if fn, ok := n.(*east.Footnote); ok && entering {
			label := string(fn.Ref)
			if _, dup := defs[label]; !dup {
				defOrder = append(defOrder, label)
			}
			defs[label] = blockText(fn, src)
			return gast.WalkSkipChildren, nil
		}
		return gast.WalkContinue, nil
	})

	lineStarts := definitionLineStarts(body)
	referenced := map[string]bool{}
	for _, m := range footnoteRef.FindAllStringSubmatchIndex(body, -1) {
		start, label := m[0], body[m[2]:m[3]]
		if inRanges(code, start) {
			continue
		}
		// "[^label]:" at the start of a line is a definition, not a reference.
		if lineStarts[start] && m[1] < len(body) && body[m[1]] == ':' {
			continue
		}
		if referenced[label] {
			continue
		}
		referenced[label] = true
		def, ok := defs[label]
		if !ok {
			warnings = append(warnings, fmt.Sprintf("footnote [^%s] is referenced but never defined", label))
			continue
		}
		footnotes = append(footnotes, Note{
			ID:        "fn-" + label,
			Label:     label,
			Number:    len(footnotes) + 1,
			Text:      def,
			RefOffset: utf8.RuneCountInString(body[:start]),
		})
	}
	// goldmark drops unreferenced definitions from its AST, so find every
	// defined label in the source to report the unused ones.
	for _, m := range footnoteRef.FindAllStringSubmatchIndex(body, -1) {
		start, label := m[0], body[m[2]:m[3]]
		if lineStarts[start] && m[1] < len(body) && body[m[1]] == ':' && !inRanges(code, start) {
			if _, seen := defs[label]; !seen {
				defs[label] = ""
				defOrder = append(defOrder, label)
			}
		}
	}
	for _, label := range defOrder {
		if !referenced[label] {
			warnings = append(warnings, fmt.Sprintf("footnote [^%s] is defined but never referenced", label))
		}
	}

	for _, span := range marginNoteSpans(body, code) {
		if span.end < 0 {
			warnings = append(warnings, fmt.Sprintf("unclosed margin note at line %d (no matching <<})", strings.Count(body[:span.start], "\n")+1))
			continue
		}
		marginNotes = append(marginNotes, Note{
			ID:        fmt.Sprintf("m%d", len(marginNotes)+1),
			Text:      strings.TrimSpace(body[span.start+3 : span.end-3]),
			RefOffset: utf8.RuneCountInString(body[:span.start]),
		})
	}
	return footnotes, marginNotes, warnings
}

// MarkdownBody is the article body for the Markdown and terminal views:
// each margin note becomes a footnote reference [^mN] with its definition
// appended, labels skipping any the author already uses, so the file reads
// as plain Markdown and no note text is lost (FR-019; research.md R6).
func MarkdownBody(body string, marginNotes []Note) string {
	if len(marginNotes) == 0 {
		return body
	}
	used := map[string]bool{}
	for _, m := range footnoteRef.FindAllStringSubmatch(body, -1) {
		used[m[1]] = true
	}
	src := []byte(body)
	code := codeRanges(notesMarkdown.Parser().Parse(text.NewReader(src)), src)

	var b strings.Builder
	var defs []string
	last, n := 0, 0
	for _, span := range marginNoteSpans(body, code) {
		if span.end < 0 {
			continue
		}
		label := ""
		for label == "" || used[label] {
			n++
			label = fmt.Sprintf("m%d", n)
		}
		used[label] = true
		b.WriteString(body[last:span.start])
		b.WriteString("[^" + label + "]")
		defs = append(defs, fmt.Sprintf("[^%s]: %s", label, strings.TrimSpace(body[span.start+3:span.end-3])))
		last = span.end
	}
	b.WriteString(body[last:])
	out := strings.TrimRight(b.String(), "\n")
	return out + "\n\n" + strings.Join(defs, "\n") + "\n"
}

// span is a byte range [start, end) of a {>> … <<} margin note; end is -1
// when the note is never closed.
type span struct{ start, end int }

// marginNoteSpans finds {>> … <<} markers outside code. A note may not
// contain another opener; an opener with no closer is reported unclosed.
func marginNoteSpans(body string, code [][2]int) []span {
	var spans []span
	for i := 0; i < len(body); {
		j := strings.Index(body[i:], "{>>")
		if j < 0 {
			break
		}
		start := i + j
		if inRanges(code, start) {
			i = start + 3
			continue
		}
		k := strings.Index(body[start+3:], "<<}")
		next := strings.Index(body[start+3:], "{>>")
		if k < 0 || (next >= 0 && next < k) {
			spans = append(spans, span{start, -1})
			i = start + 3
			continue
		}
		end := start + 3 + k + 3
		spans = append(spans, span{start, end})
		i = end
	}
	return spans
}

// codeRanges returns the byte ranges of every code span and code block in
// doc, merged and sorted.
func codeRanges(doc gast.Node, src []byte) [][2]int {
	var ranges [][2]int
	_ = gast.Walk(doc, func(n gast.Node, entering bool) (gast.WalkStatus, error) {
		if !entering {
			return gast.WalkContinue, nil
		}
		switch n.Kind() {
		case gast.KindCodeSpan:
			for c := n.FirstChild(); c != nil; c = c.NextSibling() {
				if t, ok := c.(*gast.Text); ok {
					// Widen to the surrounding backticks.
					start, stop := t.Segment.Start, t.Segment.Stop
					for start > 0 && src[start-1] == '`' {
						start--
					}
					for stop < len(src) && src[stop] == '`' {
						stop++
					}
					ranges = append(ranges, [2]int{start, stop})
				}
			}
			return gast.WalkSkipChildren, nil
		case gast.KindFencedCodeBlock, gast.KindCodeBlock:
			lines := n.Lines()
			for i := 0; i < lines.Len(); i++ {
				seg := lines.At(i)
				ranges = append(ranges, [2]int{seg.Start, seg.Stop})
			}
		}
		return gast.WalkContinue, nil
	})
	sort.Slice(ranges, func(i, j int) bool { return ranges[i][0] < ranges[j][0] })
	return ranges
}

func inRanges(ranges [][2]int, pos int) bool {
	for _, r := range ranges {
		if pos >= r[0] && pos < r[1] {
			return true
		}
	}
	return false
}

// definitionLineStarts marks every byte offset that begins a line.
func definitionLineStarts(body string) map[int]bool {
	starts := map[int]bool{0: true}
	for i := 0; i < len(body); i++ {
		if body[i] == '\n' {
			starts[i+1] = true
		}
	}
	return starts
}

// blockText is the raw Markdown of a footnote definition's paragraphs,
// joined with single spaces.
func blockText(fn gast.Node, src []byte) string {
	var parts []string
	for c := fn.FirstChild(); c != nil; c = c.NextSibling() {
		lines := c.Lines()
		var b strings.Builder
		for i := 0; i < lines.Len(); i++ {
			seg := lines.At(i)
			b.Write(seg.Value(src))
		}
		if s := strings.Join(strings.Fields(b.String()), " "); s != "" {
			parts = append(parts, s)
		}
	}
	return strings.Join(parts, " ")
}

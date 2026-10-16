// Copyright 2026 Mona Maret
// SPDX-License-Identifier: Apache-2.0

package content_test

import (
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/monamaret/biblio/internal/content"
)

// codePointOffset is the expected ref_offset of the first occurrence of
// marker in body: a Unicode code point offset, not a byte offset
// (specs/004-quarterly-site-type research.md R2).
func codePointOffset(t *testing.T, body, marker string) int {
	t.Helper()
	i := strings.Index(body, marker)
	if i < 0 {
		t.Fatalf("marker %q not in body", marker)
	}
	return utf8.RuneCountInString(body[:i])
}

// Footnotes come out in reference order with label, display number,
// inline-Markdown text and a code-point reference offset; margin notes are
// numbered m1, m2, … with their text and offset (data-model.md "Note").
func TestExtractNotes_FootnotesAndMarginNotes(t *testing.T) {
	body := "Kůň úpěl ódy.[^b] An emoji 📚 then a note.[^a]{>> cf. *Gerstner*, [1964](https://example.com) <<} End.{>>  second  <<}\n\n" +
		"[^a]: First *defined*, second referenced.\n" +
		"[^b]: Brno, 1971.\n"

	footnotes, margin, warnings := content.ExtractNotes(body)
	if len(warnings) != 0 {
		t.Errorf("warnings = %v, want none", warnings)
	}

	if len(footnotes) != 2 {
		t.Fatalf("footnotes = %+v, want 2", footnotes)
	}
	want0 := content.Note{ID: "fn-b", Label: "b", Number: 1, Text: "Brno, 1971.", RefOffset: codePointOffset(t, body, "[^b]")}
	want1 := content.Note{ID: "fn-a", Label: "a", Number: 2, Text: "First *defined*, second referenced.", RefOffset: codePointOffset(t, body, "[^a]")}
	if footnotes[0] != want0 {
		t.Errorf("footnotes[0] = %+v, want %+v", footnotes[0], want0)
	}
	if footnotes[1] != want1 {
		t.Errorf("footnotes[1] = %+v, want %+v", footnotes[1], want1)
	}
	if strings.Index(body, "[^a]") == footnotes[1].RefOffset {
		t.Error("ref_offset equals the byte offset; it must count code points (the body has multi-byte characters before it)")
	}

	if len(margin) != 2 {
		t.Fatalf("margin notes = %+v, want 2", margin)
	}
	wantM1 := content.Note{ID: "m1", Text: "cf. *Gerstner*, [1964](https://example.com)", RefOffset: codePointOffset(t, body, "{>> cf.")}
	wantM2 := content.Note{ID: "m2", Text: "second", RefOffset: codePointOffset(t, body, "{>>  second")}
	if margin[0] != wantM1 {
		t.Errorf("margin[0] = %+v, want %+v", margin[0], wantM1)
	}
	if margin[1] != wantM2 {
		t.Errorf("margin[1] = %+v, want %+v", margin[1], wantM2)
	}
}

// Markers inside code spans and fenced code blocks are never notes.
func TestExtractNotes_IgnoresCode(t *testing.T) {
	body := "Real.{>> yes <<} Code: `{>> no <<}` and `[^x]`.\n\n```\n{>> also no <<}\n[^y]\n```\n"
	footnotes, margin, warnings := content.ExtractNotes(body)
	if len(footnotes) != 0 {
		t.Errorf("footnotes = %+v, want none", footnotes)
	}
	if len(margin) != 1 || margin[0].Text != "yes" {
		t.Errorf("margin = %+v, want only the real note", margin)
	}
	if len(warnings) != 0 {
		t.Errorf("warnings = %v, want none (code is not a reference)", warnings)
	}
}

// Problems are warnings, never errors: an unclosed {>> is left as text, a
// reference with no definition and an unused definition are each reported.
func TestExtractNotes_Warnings(t *testing.T) {
	body := "Text.[^missing] More {>> never closed.\n\n[^unused]: Nobody points here.\n"
	footnotes, margin, warnings := content.ExtractNotes(body)
	if len(footnotes) != 0 || len(margin) != 0 {
		t.Errorf("notes = %+v / %+v, want none", footnotes, margin)
	}
	joined := strings.Join(warnings, "\n")
	for _, want := range []string{`footnote [^missing] is referenced but never defined`, `footnote [^unused] is defined but never referenced`, "unclosed margin note"} {
		if !strings.Contains(joined, want) {
			t.Errorf("warnings %q missing %q", joined, want)
		}
	}
}

// A footnote referenced twice is one note, numbered and offset at its
// first reference.
func TestExtractNotes_RepeatedReference(t *testing.T) {
	body := "One.[^1] Two.[^1]\n\n[^1]: Shared.\n"
	footnotes, _, _ := content.ExtractNotes(body)
	if len(footnotes) != 1 || footnotes[0].RefOffset != strings.Index(body, "[^1]") {
		t.Errorf("footnotes = %+v, want one note at the first reference", footnotes)
	}
}

// MarkdownBody rewrites margin notes as footnotes [^mN] with definitions
// appended; existing footnotes are untouched; no note text is lost
// (research.md R6; FR-019).
func TestMarkdownBody_RewritesMarginNotes(t *testing.T) {
	body := "A.[^1] B.{>> first *aside* <<} C.{>> second <<}\n\n[^1]: Real footnote.\n"
	_, margin, _ := content.ExtractNotes(body)
	got := content.MarkdownBody(body, margin)

	if strings.Contains(got, "{>>") || strings.Contains(got, "<<}") {
		t.Errorf("MarkdownBody still contains margin-note markers:\n%s", got)
	}
	for _, want := range []string{"A.[^1] B.[^m1] C.[^m2]", "[^1]: Real footnote.", "[^m1]: first *aside*", "[^m2]: second"} {
		if !strings.Contains(got, want) {
			t.Errorf("MarkdownBody missing %q:\n%s", want, got)
		}
	}
	// Round trip: the rewritten body has every note as a footnote.
	footnotes, rest, warnings := content.ExtractNotes(got)
	if len(footnotes) != 3 || len(rest) != 0 || len(warnings) != 0 {
		t.Errorf("rewritten body notes = %d footnotes, %d margin, warnings %v; want 3, 0, none", len(footnotes), len(rest), warnings)
	}
}

// Generated labels skip any label the author already uses.
func TestMarkdownBody_AvoidsLabelCollisions(t *testing.T) {
	body := "A.[^m1] B.{>> aside <<}\n\n[^m1]: Author's own m1.\n"
	_, margin, _ := content.ExtractNotes(body)
	got := content.MarkdownBody(body, margin)
	if !strings.Contains(got, "B.[^m2]") || !strings.Contains(got, "[^m2]: aside") {
		t.Errorf("MarkdownBody did not skip the used label m1:\n%s", got)
	}
	if strings.Count(got, "[^m1]:") != 1 {
		t.Errorf("label m1 defined more than once:\n%s", got)
	}
}

// A body with no margin notes is returned unchanged.
func TestMarkdownBody_NoMarginNotes(t *testing.T) {
	body := "Plain.[^1]\n\n[^1]: Note.\n"
	if got := content.MarkdownBody(body, nil); got != body {
		t.Errorf("MarkdownBody changed a body with no margin notes:\n%s", got)
	}
}

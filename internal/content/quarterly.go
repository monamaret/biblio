// Copyright 2026 Mona Maret
// SPDX-License-Identifier: Apache-2.0

package content

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode"

	"go.yaml.in/yaml/v3"
	"golang.org/x/text/unicode/norm"
)

// Quarterly file kinds (specs/004-quarterly-site-type data-model.md
// "Entities").
const (
	KindIssue    = "issue"
	KindArticle  = "article"
	KindColophon = "colophon"
)

// ParseItemKind parses and validates one content file. For quarterly sites
// the schema is chosen by the file's location — issues/<dir>/issue.md is an
// issue, any other .md directly inside an issue folder is an article, and
// colophon.md is the colophon; anything else is an error naming the file
// (FR-005–FR-008). For every other site type it is exactly ParseItem.
func ParseItemKind(relPath string, data []byte, siteType string) (Item, error) {
	if siteType != "quarterly" {
		return ParseItem(relPath, data, siteType)
	}

	kind, err := quarterlyKind(relPath)
	if err != nil {
		return Item{}, err
	}

	frontRaw, body, err := splitFrontmatter(data)
	if err != nil {
		return Item{}, fmt.Errorf("content: %s: %w", relPath, err)
	}
	fields := map[string]any{}
	if len(frontRaw) > 0 {
		if err := yaml.Unmarshal(frontRaw, &fields); err != nil {
			return Item{}, fmt.Errorf("content: %s: parsing frontmatter: %w", relPath, err)
		}
		if fields == nil {
			fields = map[string]any{}
		}
	}

	v := quarterlyValidator{path: relPath, kind: kind, fields: fields}
	switch kind {
	case KindIssue:
		v.requireString("title")
		v.requireInt("number", 1)
		v.requireDate("date")
		v.optionalString("season")
		v.optionalPair("cover", "cover_alt")
		v.optionalBool("draft")
	case KindArticle:
		v.requireString("title")
		v.requireString("category")
		v.optionalString("author")
		v.optionalString("dek")
		v.optionalString("image_caption")
		v.optionalPair("image", "image_alt")
		v.optionalInt("order")
		v.optionalStringList("tags")
	case KindColophon:
		v.requireString("title")
	}
	if v.err != nil {
		return Item{}, v.err
	}

	return Item{
		Path:   relPath,
		Slug:   SlugFromPath(relPath),
		Fields: fields,
		Body:   body,
		Kind:   kind,
	}, nil
}

// quarterlyKind maps a content-relative path to its file kind.
func quarterlyKind(relPath string) (string, error) {
	if relPath == "colophon.md" {
		return KindColophon, nil
	}
	parts := strings.Split(relPath, "/")
	if len(parts) == 3 && parts[0] == "issues" && parts[1] != "" {
		if parts[2] == "issue.md" {
			return KindIssue, nil
		}
		return KindArticle, nil
	}
	return "", fmt.Errorf("content: %s: unexpected location for a quarterly content file (want issues/<folder>/issue.md, issues/<folder>/<article>.md, or colophon.md)", relPath)
}

// IssueFolder returns the issue folder name ("01-spring-2026") for an issue
// or article path, or "" for any other path.
func IssueFolder(relPath string) string {
	parts := strings.Split(relPath, "/")
	if len(parts) == 3 && parts[0] == "issues" {
		return parts[1]
	}
	return ""
}

// quarterlyValidator checks one file's fields, keeping the first error.
// Each message names the file and the field (FR-008).
type quarterlyValidator struct {
	path   string
	kind   string
	fields map[string]any
	err    error
}

func (v *quarterlyValidator) fail(field, problem string) {
	if v.err == nil {
		v.err = fmt.Errorf("content: %s: field %q %s", v.path, field, problem)
	}
}

func (v *quarterlyValidator) missing(field string) {
	if v.err == nil {
		v.err = fmt.Errorf("content: %s: missing required frontmatter field %q (required for an %s)", v.path, field, kindNoun(v.kind))
	}
}

func kindNoun(kind string) string {
	if kind == KindColophon {
		return "colophon page"
	}
	return kind
}

func (v *quarterlyValidator) present(field string) bool {
	val, ok := v.fields[field]
	return ok && val != nil
}

func (v *quarterlyValidator) requireString(field string) {
	if !v.present(field) {
		v.missing(field)
		return
	}
	s, ok := v.fields[field].(string)
	if !ok {
		v.fail(field, "must be a string")
		return
	}
	if strings.TrimSpace(s) == "" {
		v.missing(field)
	}
}

func (v *quarterlyValidator) optionalString(field string) {
	if !v.present(field) {
		return
	}
	s, ok := v.fields[field].(string)
	if !ok || strings.TrimSpace(s) == "" {
		v.fail(field, "must be a non-empty string when present")
	}
}

func (v *quarterlyValidator) requireInt(field string, min int) {
	if !v.present(field) {
		v.missing(field)
		return
	}
	v.checkInt(field, min)
}

func (v *quarterlyValidator) optionalInt(field string) {
	if v.present(field) {
		v.checkInt(field, -1<<31)
	}
}

func (v *quarterlyValidator) checkInt(field string, min int) {
	n, ok := v.fields[field].(int)
	if !ok {
		v.fail(field, "must be an integer")
		return
	}
	if n < min {
		v.fail(field, fmt.Sprintf("must be %d or more", min))
	}
}

// requireDate accepts a quoted "YYYY-MM-DD" string or an unquoted YAML
// date, and normalizes the field to the YYYY-MM-DD string.
func (v *quarterlyValidator) requireDate(field string) {
	if !v.present(field) {
		v.missing(field)
		return
	}
	switch d := v.fields[field].(type) {
	case time.Time:
		v.fields[field] = d.Format("2006-01-02")
	case string:
		if _, err := time.Parse("2006-01-02", d); err != nil {
			v.fail(field, "is not a valid date (want YYYY-MM-DD)")
		}
	default:
		v.fail(field, "is not a valid date (want YYYY-MM-DD)")
	}
}

func (v *quarterlyValidator) optionalBool(field string) {
	if !v.present(field) {
		return
	}
	if _, ok := v.fields[field].(bool); !ok {
		v.fail(field, "must be true or false")
	}
}

// optionalPair validates an optional image path and its required alt text.
func (v *quarterlyValidator) optionalPair(field, altField string) {
	if !v.present(field) {
		return
	}
	if s, ok := v.fields[field].(string); !ok || strings.TrimSpace(s) == "" {
		v.fail(field, "must be a non-empty path when present")
		return
	}
	if alt, ok := v.fields[altField].(string); !ok || strings.TrimSpace(alt) == "" {
		if v.err == nil {
			v.err = fmt.Errorf("content: %s: field %q is required when %q is set", v.path, altField, field)
		}
	}
}

func (v *quarterlyValidator) optionalStringList(field string) {
	if !v.present(field) {
		return
	}
	list, ok := v.fields[field].([]any)
	if !ok {
		v.fail(field, "must be a list of strings")
		return
	}
	for _, el := range list {
		if _, ok := el.(string); !ok {
			v.fail(field, fmt.Sprintf("contains a non-string element %v", el))
			return
		}
	}
}

// asciiFold maps letters NFKD doesn't decompose to a base letter.
var asciiFold = map[rune]string{
	'ł': "l", 'Ł': "l", 'đ': "d", 'Đ': "d", 'ø': "o", 'Ø': "o",
	'ß': "ss", 'æ': "ae", 'Æ': "ae", 'œ': "oe", 'Œ': "oe", 'ı': "i",
}

// CategoryKey normalizes a free-text category label into the key used for
// grouping and for its /category/<key> URL: diacritics folded to ASCII,
// lowercased, runs of anything else collapsed to one hyphen. Labels that
// differ only in case or spacing share a key (FR-011; research.md R10).
func CategoryKey(label string) (string, error) {
	var b strings.Builder
	pendingHyphen := false
	for _, r := range norm.NFKD.String(label) {
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		var piece string
		switch {
		case asciiFold[r] != "":
			piece = asciiFold[r]
		case r < unicode.MaxASCII && (unicode.IsLetter(r) || unicode.IsDigit(r)):
			piece = strings.ToLower(string(r))
		default:
			pendingHyphen = b.Len() > 0
			continue
		}
		if pendingHyphen {
			b.WriteByte('-')
			pendingHyphen = false
		}
		b.WriteString(piece)
	}
	if b.Len() == 0 {
		return "", errors.New("category label has no letters or digits")
	}
	return b.String(), nil
}

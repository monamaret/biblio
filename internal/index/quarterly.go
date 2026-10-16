// Copyright 2026 Mona Maret
// SPDX-License-Identifier: Apache-2.0

package index

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/monamaret/biblio/internal/content"
)

// SiteMeta is the index's `site` object: the site.yaml name and
// description, for page titles, the quarterly masthead and llms.txt.
type SiteMeta struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// MarshalJSON keeps profile, log, and docs indexes exactly
// {"site": {...}, "entries": [...]}, and always emits every quarterly key — empty arrays
// and a null current_issue included — for quarterly indexes, which are the
// ones carrying Typesetting (specs/004-quarterly-site-type data-model.md
// "content-index.json shape").
func (idx Index) MarshalJSON() ([]byte, error) {
	if idx.Typesetting == nil {
		return json.Marshal(struct {
			Site    *SiteMeta `json:"site,omitempty"`
			Entries []Entry   `json:"entries"`
		}{idx.Site, idx.Entries})
	}
	issues, categories := idx.Issues, idx.Categories
	if issues == nil {
		issues = []IssueSummary{}
	}
	if categories == nil {
		categories = []CategorySummary{}
	}
	return json.Marshal(struct {
		Site         *SiteMeta         `json:"site"`
		Entries      []Entry           `json:"entries"`
		Issues       []IssueSummary    `json:"issues"`
		Categories   []CategorySummary `json:"categories"`
		CurrentIssue *string           `json:"current_issue"`
		Typesetting  *IndexTypesetting `json:"typesetting"`
	}{idx.Site, idx.Entries, issues, categories, idx.CurrentIssue, idx.Typesetting})
}

// quarterlyIssue gathers one issue folder's parsed files.
type quarterlyIssue struct {
	folder   string
	item     *content.Item
	articles []content.Item
}

func (q *quarterlyIssue) number() int { return q.item.Fields["number"].(int) }
func (q *quarterlyIssue) draft() bool { d, _ := q.item.Fields["draft"].(bool); return d }

// generateQuarterly builds a quarterly index: every file validated by its
// kind, cross-file rules checked, drafts handled per opts, and the derived
// fields and site-level arrays added (data-model.md "Ordering rules",
// "Drafts", "Validation summary"). All errors are collected.
func generateQuarterly(contentDir string, cfg content.SiteConfig, opts Options) (Index, []error) {
	var errs []error
	issues := map[string]*quarterlyIssue{}
	var colophon *content.Item

	issueFor := func(folder string) *quarterlyIssue {
		if issues[folder] == nil {
			issues[folder] = &quarterlyIssue{folder: folder}
		}
		return issues[folder]
	}

	walkErr := filepath.WalkDir(contentDir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(p, ".md") {
			return nil
		}
		rel, err := filepath.Rel(contentDir, p)
		if err != nil {
			errs = append(errs, fmt.Errorf("index: resolving path %q: %w", p, err))
			return nil
		}
		rel = filepath.ToSlash(rel)
		data, err := os.ReadFile(p)
		if err != nil {
			errs = append(errs, fmt.Errorf("index: reading %q: %w", rel, err))
			return nil
		}
		item, err := content.ParseItemKind(rel, data, "quarterly")
		if err != nil {
			errs = append(errs, err)
			return nil
		}
		switch item.Kind {
		case content.KindColophon:
			colophon = &item
		case content.KindIssue:
			issueFor(content.IssueFolder(rel)).item = &item
		case content.KindArticle:
			q := issueFor(content.IssueFolder(rel))
			q.articles = append(q.articles, item)
		}
		return nil
	})
	if walkErr != nil {
		errs = append(errs, fmt.Errorf("index: walking content directory: %w", walkErr))
	}

	folders := make([]string, 0, len(issues))
	for f := range issues {
		folders = append(folders, f)
	}
	sort.Strings(folders)

	// Cross-file rules: every issue folder has an issue.md; numbers are
	// unique across all issues, drafts included.
	byNumber := map[int]string{}
	var all []*quarterlyIssue
	for _, f := range folders {
		q := issues[f]
		if q.item == nil {
			errs = append(errs, fmt.Errorf("content: issues/%s/: issue folder %q has articles but no issue.md", f, f))
			continue
		}
		if other, dup := byNumber[q.number()]; dup {
			errs = append(errs, fmt.Errorf("content: issues %q and %q both have number %d (issue numbers must be unique)", other, f, q.number()))
			continue
		}
		byNumber[q.number()] = f
		all = append(all, q)
	}
	if len(errs) > 0 {
		return Index{}, errs
	}

	idx := Index{
		Site:        &SiteMeta{Name: cfg.Name, Description: cfg.Description},
		Typesetting: opts.Typesetting,
	}
	if idx.Typesetting == nil {
		idx.Typesetting = &IndexTypesetting{MarginNotes: true, Numerals: "roman"}
	}

	var included []*quarterlyIssue
	for _, q := range all {
		if q.draft() && !opts.IncludeDrafts {
			idx.DraftIssuesExcluded++
			continue
		}
		included = append(included, q)
	}
	sort.SliceStable(included, func(i, j int) bool { return included[i].number() > included[j].number() })

	catIndex := map[string]int{}
	for _, q := range included {
		ordered := contentsOrder(q.articles)
		draft := q.draft()
		issueFields := copyFields(q.item.Fields)
		issueFields["kind"] = content.KindIssue
		issueFields["issue"] = q.folder
		issueFields["article_count"] = len(ordered)
		coverPlate := plateFor(opts.Plates, q.item.Path, issueFields["cover"])
		if coverPlate != "" {
			issueFields["cover_plate"] = coverPlate
		}
		idx.Entries = append(idx.Entries, Entry{
			Slug:   "issues/" + q.folder,
			Path:   q.item.Path,
			Body:   q.item.Body,
			Fields: issueFields,
		})

		summary := IssueSummary{
			Issue:        q.folder,
			Slug:         "issues/" + q.folder,
			Number:       q.number(),
			Title:        q.item.Fields["title"].(string),
			Date:         q.item.Fields["date"].(string),
			CoverPlate:   coverPlate,
			Draft:        draft,
			ArticleCount: len(ordered),
		}
		summary.Season, _ = q.item.Fields["season"].(string)
		summary.CoverAlt, _ = q.item.Fields["cover_alt"].(string)
		idx.Issues = append(idx.Issues, summary)
		if !draft && idx.CurrentIssue == nil {
			cur := q.folder
			idx.CurrentIssue = &cur
		}

		for pos, a := range ordered {
			label := strings.TrimSpace(a.Fields["category"].(string))
			key, err := content.CategoryKey(label)
			if err != nil {
				errs = append(errs, fmt.Errorf("content: %s: field %q: %v", a.Path, "category", err))
				continue
			}
			fields := copyFields(a.Fields)
			fields["kind"] = content.KindArticle
			fields["issue"] = q.folder
			fields["issue_number"] = q.number()
			fields["issue_date"] = q.item.Fields["date"]
			fields["category_key"] = key
			fields["position"] = pos + 1
			if draft {
				fields["draft"] = true
			}
			if plate := plateFor(opts.Plates, a.Path, fields["image"]); plate != "" {
				fields["plate"] = plate
			}
			body := content.RewriteImages(a.Body, path.Dir(a.Path), opts.Plates)
			footnotes, marginNotes, warnings := content.ExtractNotes(body)
			if footnotes == nil {
				footnotes = []content.Note{}
			}
			if marginNotes == nil {
				marginNotes = []content.Note{}
			}
			fields["footnotes"] = footnotes
			fields["margin_notes"] = marginNotes
			fields["markdown_body"] = content.MarkdownBody(body, marginNotes)
			for _, w := range warnings {
				idx.Warnings = append(idx.Warnings, a.Path+": "+w)
			}
			idx.Entries = append(idx.Entries, Entry{Slug: a.Slug, Path: a.Path, Body: body, Fields: fields})

			if i, seen := catIndex[key]; seen {
				idx.Categories[i].ArticleCount++
			} else {
				catIndex[key] = len(idx.Categories)
				idx.Categories = append(idx.Categories, CategorySummary{Key: key, Label: label, ArticleCount: 1})
			}
		}
	}

	if colophon != nil {
		fields := copyFields(colophon.Fields)
		fields["kind"] = content.KindColophon
		idx.Entries = append(idx.Entries, Entry{Slug: colophon.Slug, Path: colophon.Path, Body: colophon.Body, Fields: fields})
	}

	if len(errs) > 0 {
		return Index{}, errs
	}
	return idx, nil
}

// contentsOrder sorts an issue's articles for its contents page: by
// `order` ascending (unordered articles after ordered ones), ties by file
// name; then grouped by category key, categories in order of their first
// article (data-model.md "Ordering rules").
func contentsOrder(articles []content.Item) []content.Item {
	sorted := append([]content.Item(nil), articles...)
	sort.SliceStable(sorted, func(i, j int) bool {
		oi, iok := sorted[i].Fields["order"].(int)
		oj, jok := sorted[j].Fields["order"].(int)
		switch {
		case iok && jok && oi != oj:
			return oi < oj
		case iok != jok:
			return iok
		}
		return path.Base(sorted[i].Path) < path.Base(sorted[j].Path)
	})

	var keys []string
	groups := map[string][]content.Item{}
	for _, a := range sorted {
		key, err := content.CategoryKey(a.Fields["category"].(string))
		if err != nil {
			key = ""
		}
		if _, seen := groups[key]; !seen {
			keys = append(keys, key)
		}
		groups[key] = append(groups[key], a)
	}
	var out []content.Item
	for _, k := range keys {
		out = append(out, groups[k]...)
	}
	return out
}

func copyFields(fields map[string]any) map[string]any {
	out := make(map[string]any, len(fields)+8)
	for k, v := range fields {
		out[k] = v
	}
	return out
}

// plateFor resolves a frontmatter image path (relative to the file at
// relPath) against the plate map, returning "" when there's no plate.
func plateFor(plates map[string]string, relPath string, image any) string {
	src, ok := image.(string)
	if !ok || plates == nil {
		return ""
	}
	return plates[path.Join(path.Dir(relPath), src)]
}

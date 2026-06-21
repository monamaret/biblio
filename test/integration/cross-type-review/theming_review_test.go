// Package crosstypereview is the structured, repeatable comparison run only
// once all three site types exist (research.md "Why documentation is
// replicated last"; data-model.md "Cross-Type Consistency Record"). This
// file covers the theming half (FR-003, Scenario 2, SC-002).
package crosstypereview_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/monamaret/biblio/internal/scaffold"
	"github.com/monamaret/biblio/internal/theme"
)

// scaffoldTokensCSS scaffolds siteType and returns its generated
// tokens.css content (data-model.md "Cross-Type Consistency Record"
// theming-check).
func scaffoldTokensCSS(t *testing.T, siteType string) string {
	t.Helper()
	presets := theme.PresetsForType(siteType)
	if len(presets) == 0 {
		t.Fatalf("theme.PresetsForType(%q): no presets defined", siteType)
	}

	destDir := t.TempDir()
	projectDir, err := scaffold.Generate(destDir, scaffold.Options{
		ProjectName: "review-" + siteType,
		SiteType:    siteType,
		Domain:      "example.com",
		Theme:       presets[0].ID,
		MCPEnabled:  true,
	})
	if err != nil {
		t.Fatalf("scaffold.Generate(%q): %v", siteType, err)
	}

	data, err := os.ReadFile(filepath.Join(projectDir, "public", "tokens.css"))
	if err != nil {
		t.Fatalf("reading public/tokens.css for %q: %v", siteType, err)
	}
	return string(data)
}

// TestCrossTypeTheming_EachSiteRendersItsOwnPreset is the theming half of
// the Cross-Type Consistency Record: each site type's generated tokens.css
// must reflect exactly its own preset's token values — never a value
// leaked or bled through from another type's preset (FR-003, Scenario 2,
// SC-002). This is the correctness bar: it is *not* "every token value
// must differ across types" — presets intentionally share some values
// (e.g. profile/log/docs all use the same `font-heading`), and that
// overlap is a legitimate design choice, not bleed-through.
func TestCrossTypeTheming_EachSiteRendersItsOwnPreset(t *testing.T) {
	siteTypes := []string{"profile", "log", "docs"}

	for _, siteType := range siteTypes {
		siteType := siteType
		t.Run(siteType, func(t *testing.T) {
			css := scaffoldTokensCSS(t, siteType)
			presets := theme.PresetsForType(siteType)
			preset := presets[0]

			for token, want := range preset.Tokens {
				wantLine := "--" + token + ": " + want + ";"
				if !strings.Contains(css, wantLine) {
					t.Errorf("tokens.css for %q: missing or mismatched %q (want value %q); got:\n%s", siteType, token, want, css)
				}
			}
		})
	}
}

// TestCrossTypeTheming_NoBleedThrough confirms each site type's distinguishing
// tokens (color-fg, color-accent — the tokens each preset defines with a
// genuinely distinct value per type) never match another type's value for
// that same token, which would indicate a generation bug routing the wrong
// preset into the wrong site (spec Edge Cases — "shared token resolving
// differently per type" is the failure mode this guards against, inverted:
// here we guard against a *site* resolving a *different* type's token).
//
// color-bg is deliberately excluded from this distinguishing set: profile
// and docs both define a white (#ffffff) background — two light themes
// alongside log's dark one — which is a legitimate, intentional design
// choice (the same category as all three types sharing one font-heading
// value), not bleed-through. Asserting it must differ here would be a false
// positive against a real design decision, not a defect.
func TestCrossTypeTheming_NoBleedThrough(t *testing.T) {
	css := map[string]string{
		"profile": scaffoldTokensCSS(t, "profile"),
		"log":     scaffoldTokensCSS(t, "log"),
		"docs":    scaffoldTokensCSS(t, "docs"),
	}

	distinguishing := []string{"color-fg", "color-accent"}
	for _, token := range distinguishing {
		seen := map[string]string{} // value -> first site type that had it
		for _, siteType := range []string{"profile", "log", "docs"} {
			preset := theme.PresetsForType(siteType)[0]
			want := preset.Tokens[token]
			if other, ok := seen[want]; ok {
				t.Errorf("token %q: %q and %q both resolve to %q — expected these distinguishing tokens to differ per type", token, other, siteType, want)
			}
			seen[want] = siteType

			wantLine := "--" + token + ": " + want + ";"
			if !strings.Contains(css[siteType], wantLine) {
				t.Errorf("tokens.css for %q: token %q does not match its own preset (want %q); got:\n%s", siteType, token, want, css[siteType])
			}
		}
	}
}

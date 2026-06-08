package theme_test

import (
	"strings"
	"testing"

	"github.com/monamaret/biblio/internal/theme"
)

// One ready-to-use preset per site type (FR-004; data-model.md "Theme Preset").
func TestPresetsForType_OnePerSiteType(t *testing.T) {
	for _, siteType := range []string{"profile", "log", "docs"} {
		presets := theme.PresetsForType(siteType)
		if len(presets) == 0 {
			t.Errorf("PresetsForType(%q): want at least one ready-to-use preset, got none", siteType)
			continue
		}
		for _, p := range presets {
			if p.SiteType != siteType {
				t.Errorf("PresetsForType(%q): preset %q has SiteType %q, want %q", siteType, p.ID, p.SiteType, siteType)
			}
			if p.ID == "" {
				t.Errorf("PresetsForType(%q): preset has an empty ID", siteType)
			}
		}
	}
}

// tokens.css generation must be deterministic and owned solely by biblio-cli
// (data-model.md "Theme Preset" validation rules).
func TestGenerateTokensCSS_Deterministic(t *testing.T) {
	presets := theme.PresetsForType("profile")
	if len(presets) == 0 {
		t.Fatal(`PresetsForType("profile"): want at least one preset to generate from`)
	}
	preset := presets[0]

	first, err := theme.GenerateTokensCSS(preset)
	if err != nil {
		t.Fatalf("GenerateTokensCSS() first call: %v", err)
	}
	second, err := theme.GenerateTokensCSS(preset)
	if err != nil {
		t.Fatalf("GenerateTokensCSS() second call: %v", err)
	}
	if first != second {
		t.Errorf("GenerateTokensCSS() is not deterministic for the same preset:\nfirst:\n%s\nsecond:\n%s", first, second)
	}

	if !strings.Contains(first, ":root") {
		t.Errorf("GenerateTokensCSS() output missing a :root custom-property block:\n%s", first)
	}
	for name := range preset.Tokens {
		if !strings.Contains(first, "--"+name) {
			t.Errorf("GenerateTokensCSS() output missing custom property for token %q:\n%s", name, first)
		}
	}
}

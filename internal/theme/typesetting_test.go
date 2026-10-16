// Copyright 2026 Mona Maret
// SPDX-License-Identifier: Apache-2.0

package theme_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/monamaret/biblio/internal/theme"
)

// Every typesetting key accepts its documented values and rejects anything
// else with a message naming the key and the accepted values (FR-023;
// specs/004-quarterly-site-type contracts/config-and-tokens.md "Error
// format"; research.md R9).
func TestParseTypesetting_ValidValues(t *testing.T) {
	valid := map[string]any{
		"accent":           "#1e4fd7",
		"font_body":        "Söhne",
		"font_display":     "Archivo Black",
		"font_mono":        "IBM Plex Mono",
		"base_size":        "1rem",
		"scale":            1.333,
		"leading":          1.6,
		"measure":          "60ch",
		"align":            "justify",
		"grid":             "centered",
		"margin_notes":     false,
		"drop_caps":        "raised",
		"lead_in":          "none",
		"numerals":         "arabic",
		"rules":            "hairline",
		"fleurons":         true,
		"grayscale_images": false,
	}
	ts, errs := theme.ParseTypesetting(valid)
	if len(errs) != 0 {
		t.Fatalf("ParseTypesetting(valid) errors = %v, want none", errs)
	}
	if ts.Accent == nil || *ts.Accent != "#1e4fd7" {
		t.Errorf("Accent = %v, want #1e4fd7", ts.Accent)
	}
	if ts.Scale == nil || *ts.Scale != 1.333 {
		t.Errorf("Scale = %v, want 1.333", ts.Scale)
	}
	if ts.MarginNotes == nil || *ts.MarginNotes {
		t.Errorf("MarginNotes = %v, want false", ts.MarginNotes)
	}
	// Integers are valid numbers for scale/leading (YAML `leading: 2`).
	if _, errs := theme.ParseTypesetting(map[string]any{"leading": 2}); len(errs) != 0 {
		t.Errorf("ParseTypesetting(leading: 2) errors = %v, want none", errs)
	}
	// base_size in px is accepted within range.
	if _, errs := theme.ParseTypesetting(map[string]any{"base_size": "17px"}); len(errs) != 0 {
		t.Errorf("ParseTypesetting(base_size: 17px) errors = %v, want none", errs)
	}
}

func TestParseTypesetting_InvalidValues(t *testing.T) {
	cases := []struct {
		key   string
		value any
		want  []string // substrings the single error must contain
	}{
		{"accent", "red", []string{"typesetting.accent", "#rrggbb"}},
		{"accent", "#12345", []string{"typesetting.accent"}},
		{"font_body", "x;}", []string{"typesetting.font_body", "letters"}},
		{"font_display", `Inter", serif`, []string{"typesetting.font_display"}},
		{"font_mono", "", []string{"typesetting.font_mono"}},
		{"base_size", "3rem", []string{"typesetting.base_size", "0.75rem"}},
		{"base_size", "1em", []string{"typesetting.base_size", "rem", "px"}},
		{"scale", 1.0, []string{"typesetting.scale", "1.1", "1.618"}},
		{"scale", 2, []string{"typesetting.scale", "out of range"}},
		{"scale", "big", []string{"typesetting.scale", "number"}},
		{"leading", 3.0, []string{"typesetting.leading", "1.1", "2"}},
		{"measure", "34px", []string{"typesetting.measure", "em", "ch"}},
		{"measure", "120em", []string{"typesetting.measure", "20", "90"}},
		{"align", "center", []string{"typesetting.align", "left", "justify"}},
		{"grid", "golden", []string{"typesetting.grid", "asymmetric", "centered"}},
		{"margin_notes", "yes", []string{"typesetting.margin_notes", "true or false"}},
		{"drop_caps", "big", []string{"typesetting.drop_caps", "none", "drop", "raised"}},
		{"lead_in", "caps", []string{"typesetting.lead_in", "none", "small_caps"}},
		{"numerals", "greek", []string{"typesetting.numerals", "roman", "arabic"}},
		{"rules", "thick", []string{"typesetting.rules", "none", "hairline", "heavy"}},
		{"fleurons", 1, []string{"typesetting.fleurons", "true or false"}},
		{"grayscale_images", "no", []string{"typesetting.grayscale_images", "true or false"}},
		{"colour", "#000000", []string{"typesetting.colour", "unknown key", "accent", "grayscale_images"}},
	}
	for _, tc := range cases {
		t.Run(tc.key, func(t *testing.T) {
			_, errs := theme.ParseTypesetting(map[string]any{tc.key: tc.value})
			if len(errs) != 1 {
				t.Fatalf("ParseTypesetting(%s: %v) errors = %v, want exactly one", tc.key, tc.value, errs)
			}
			msg := errs[0].Error()
			for _, sub := range tc.want {
				if !strings.Contains(msg, sub) {
					t.Errorf("error %q does not contain %q", msg, sub)
				}
			}
		})
	}
}

// Every problem is reported, not just the first (contracts/config-and-tokens.md
// "Error format": "One line per problem, all collected").
func TestParseTypesetting_CollectsAllErrors(t *testing.T) {
	_, errs := theme.ParseTypesetting(map[string]any{
		"scale":  2,
		"align":  "center",
		"colour": "#000000",
	})
	if len(errs) != 3 {
		t.Fatalf("ParseTypesetting errors = %v, want 3", errs)
	}
	// Deterministic order: sorted by key.
	if !strings.Contains(errs[0].Error(), "typesetting.align") ||
		!strings.Contains(errs[1].Error(), "typesetting.colour") ||
		!strings.Contains(errs[2].Error(), "typesetting.scale") {
		t.Errorf("errors not sorted by key: %v", errs)
	}
}

// With no overrides, the effective values are the quarterly-default preset
// values (research doc 004 typesetting schema defaults).
func TestResolve_NoOverridesGivesDefaults(t *testing.T) {
	preset, ok := theme.PresetByID("quarterly-default")
	if !ok {
		t.Fatal(`PresetByID("quarterly-default"): not found`)
	}
	eff, err := theme.Resolve(preset, theme.Typesetting{})
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	want := theme.Effective{
		Accent:          "#d7261e",
		FontBody:        "Archivo",
		FontDisplay:     "Archivo",
		FontMono:        "IBM Plex Mono",
		BaseSize:        "1.0625rem",
		Scale:           1.25,
		Leading:         1.5,
		Measure:         "34em",
		Align:           "left",
		Grid:            "asymmetric",
		MarginNotes:     true,
		DropCaps:        "none",
		LeadIn:          "small_caps",
		Numerals:        "roman",
		Rules:           "heavy",
		Fleurons:        false,
		GrayscaleImages: true,
	}
	if eff.Overridden != nil {
		t.Errorf("Overridden = %v, want nil with no overrides", eff.Overridden)
	}
	eff.Overridden = nil
	if !reflect.DeepEqual(eff, want) {
		t.Errorf("Resolve(no overrides) =\n%+v\nwant\n%+v", eff, want)
	}
}

// One override changes only that value and is listed in Overridden.
func TestResolve_OneOverrideChangesOnlyThatValue(t *testing.T) {
	preset, _ := theme.PresetByID("quarterly-default")
	ts, errs := theme.ParseTypesetting(map[string]any{"numerals": "arabic"})
	if len(errs) != 0 {
		t.Fatalf("ParseTypesetting errors = %v", errs)
	}
	base, _ := theme.Resolve(preset, theme.Typesetting{})
	eff, err := theme.Resolve(preset, ts)
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	if eff.Numerals != "arabic" {
		t.Errorf("Numerals = %q, want arabic", eff.Numerals)
	}
	if len(eff.Overridden) != 1 || eff.Overridden[0] != "numerals" {
		t.Errorf("Overridden = %v, want [numerals]", eff.Overridden)
	}
	eff.Numerals, eff.Overridden = base.Numerals, nil
	if !reflect.DeepEqual(eff, base) {
		t.Errorf("an override changed other values:\n%+v\nwant\n%+v", eff, base)
	}
}

// Resolve only applies to the quarterly preset — the other types have no
// typesetting schema.
func TestResolve_RejectsNonQuarterlyPreset(t *testing.T) {
	preset, _ := theme.PresetByID("log-default")
	if _, err := theme.Resolve(preset, theme.Typesetting{}); err == nil {
		t.Error("Resolve(log-default) error = nil, want an error")
	}
}

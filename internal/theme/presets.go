// Package theme owns the Theme Preset entity and the deterministic
// generation of tokens.css from a selected preset (FR-004; Principle II —
// theming logic lives only here, never duplicated on the TS side).
package theme

// Preset is a named, ready-to-use visual-theme definition selectable at
// scaffold time — exactly one per site type in this release (FR-004;
// data-model.md "Theme Preset").
type Preset struct {
	// ID is the identifier passed via --theme; unique across all presets.
	ID string
	// SiteType is the preset's matching site type (profile, log, or docs).
	SiteType string
	// Name is the preset's human-readable display name (used by the guided
	// `new` flow's theme prompt).
	Name string
	// Tokens maps design-token names to the CSS custom-property values
	// tokens.css is generated from.
	Tokens map[string]string
}

// presets holds the one ready-to-use preset shipped per site type. Adding a
// second preset for a type is a future-release concern (FR-004 only requires
// "a distinct, ready-to-use preset" per type for this release).
var presets = []Preset{
	{
		ID:       "profile-default",
		SiteType: "profile",
		Name:     "Profile Default",
		Tokens: map[string]string{
			"color-bg":     "#ffffff",
			"color-fg":     "#1a1a1a",
			"color-accent": "#2563eb",
			"color-muted":  "#6b7280",
			"font-body":    `"Inter", system-ui, sans-serif`,
			"font-heading": `"Inter", system-ui, sans-serif`,
		},
	},
	{
		ID:       "log-default",
		SiteType: "log",
		Name:     "Log Default",
		Tokens: map[string]string{
			"color-bg":     "#0f172a",
			"color-fg":     "#e2e8f0",
			"color-accent": "#38bdf8",
			"color-muted":  "#94a3b8",
			"font-body":    `"Source Serif 4", Georgia, serif`,
			"font-heading": `"Inter", system-ui, sans-serif`,
		},
	},
	{
		ID:       "docs-default",
		SiteType: "docs",
		Name:     "Docs Default",
		Tokens: map[string]string{
			"color-bg":     "#ffffff",
			"color-fg":     "#1f2937",
			"color-accent": "#7c3aed",
			"color-muted":  "#6b7280",
			"font-body":    `"Inter", system-ui, sans-serif`,
			"font-heading": `"Inter", system-ui, sans-serif`,
		},
	},
}

// PresetsForType returns the ready-to-use presets defined for siteType, in
// the order theme choices should be offered (data-model.md "Theme Preset" —
// "each site-type MUST have at least one ready-to-use preset").
func PresetsForType(siteType string) []Preset {
	var matches []Preset
	for _, p := range presets {
		if p.SiteType == siteType {
			matches = append(matches, p)
		}
	}
	return matches
}

// PresetByID resolves a preset by its --theme identifier.
func PresetByID(id string) (Preset, bool) {
	for _, p := range presets {
		if p.ID == id {
			return p, true
		}
	}
	return Preset{}, false
}

package theme

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// Typesetting holds the author's optional overrides from the quarterly
// `typesetting:` block in content/site.yaml. Every field is a pointer so
// "not set" (use the preset default) is distinct from any set value
// (specs/004-quarterly-site-type FR-022; research.md R9). A Typesetting is
// only ever produced by ParseTypesetting, so every set value has already
// passed a closed validation pattern before it can reach tokens.css.
type Typesetting struct {
	Accent          *string
	FontBody        *string
	FontDisplay     *string
	FontMono        *string
	BaseSize        *string
	Scale           *float64
	Leading         *float64
	Measure         *string
	Align           *string
	Grid            *string
	MarginNotes     *bool
	DropCaps        *string
	LeadIn          *string
	Numerals        *string
	Rules           *string
	Fleurons        *bool
	GrayscaleImages *bool
}

// Effective is the resolved typesetting for a quarterly site: the
// quarterly-default values with every set override applied. It is the one
// input both tokens.css generation and the content index's structural
// switches are produced from (research.md R3).
type Effective struct {
	Accent          string
	FontBody        string
	FontDisplay     string
	FontMono        string
	BaseSize        string
	Scale           float64
	Leading         float64
	Measure         string
	Align           string
	Grid            string
	MarginNotes     bool
	DropCaps        string
	LeadIn          string
	Numerals        string
	Rules           string
	Fleurons        bool
	GrayscaleImages bool
	// Overridden lists the typesetting keys the author set, sorted; nil
	// when nothing was overridden.
	Overridden []string
}

// quarterlyDefaults are the typesetting defaults of the quarterly-default
// preset (specs/research/004-quarterly-site-type.md "Typesetting schema").
var quarterlyDefaults = Effective{
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

var (
	colorPattern  = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)
	familyPattern = regexp.MustCompile(`^[\p{L}\p{N}][\p{L}\p{N} \-]{0,63}$`)
	lengthPattern = regexp.MustCompile(`^(\d+(?:\.\d+)?)([a-z]+)$`)
)

// typesettingKey describes how one `typesetting:` key is validated and
// applied to a Typesetting.
type typesettingKey struct {
	parse func(key string, v any, ts *Typesetting) error
}

func enumKey(field func(*Typesetting) **string, values ...string) typesettingKey {
	return typesettingKey{parse: func(key string, v any, ts *Typesetting) error {
		s, ok := v.(string)
		if !ok || !contains(values, s) {
			return fmt.Errorf("typesetting.%s: %s is not valid (want one of: %s)", key, show(v), strings.Join(values, ", "))
		}
		*field(ts) = &s
		return nil
	}}
}

func boolKey(field func(*Typesetting) **bool) typesettingKey {
	return typesettingKey{parse: func(key string, v any, ts *Typesetting) error {
		b, ok := v.(bool)
		if !ok {
			return fmt.Errorf("typesetting.%s: %s is not valid (want true or false)", key, show(v))
		}
		*field(ts) = &b
		return nil
	}}
}

func familyKey(field func(*Typesetting) **string) typesettingKey {
	return typesettingKey{parse: func(key string, v any, ts *Typesetting) error {
		s, ok := v.(string)
		if !ok || !familyPattern.MatchString(s) {
			return fmt.Errorf("typesetting.%s: %s is not a valid font family name (want letters, digits, spaces and hyphens, up to 64 characters)", key, show(v))
		}
		*field(ts) = &s
		return nil
	}}
}

func numberKey(field func(*Typesetting) **float64, min, max float64) typesettingKey {
	want := fmt.Sprintf("want a number from %s to %s", num(min), num(max))
	return typesettingKey{parse: func(key string, v any, ts *Typesetting) error {
		f, ok := toFloat(v)
		if !ok {
			return fmt.Errorf("typesetting.%s: %s is not a number (%s)", key, show(v), want)
		}
		if f < min || f > max {
			return fmt.Errorf("typesetting.%s: %s is out of range (%s)", key, num(f), want)
		}
		*field(ts) = &f
		return nil
	}}
}

// lengthKey validates a CSS length whose unit is one of units, with the
// numeric part inside that unit's [min, max] range.
func lengthKey(field func(*Typesetting) **string, example string, ranges map[string][2]float64, rangeText string) typesettingKey {
	units := make([]string, 0, len(ranges))
	for u := range ranges {
		units = append(units, u)
	}
	sort.Strings(units)
	return typesettingKey{parse: func(key string, v any, ts *Typesetting) error {
		s, _ := v.(string)
		m := lengthPattern.FindStringSubmatch(s)
		if m == nil || ranges[m[2]] == [2]float64{} {
			return fmt.Errorf("typesetting.%s: %s is not valid (want a length in %s, e.g. %q)", key, show(v), strings.Join(units, " or "), example)
		}
		f, _ := strconv.ParseFloat(m[1], 64)
		if r := ranges[m[2]]; f < r[0] || f > r[1] {
			return fmt.Errorf("typesetting.%s: %q is out of range (want %s)", key, s, rangeText)
		}
		*field(ts) = &s
		return nil
	}}
}

var typesettingKeys = map[string]typesettingKey{
	"accent": {parse: func(key string, v any, ts *Typesetting) error {
		s, ok := v.(string)
		if !ok || !colorPattern.MatchString(s) {
			return fmt.Errorf("typesetting.%s: %s is not a valid color (want #rrggbb, e.g. \"#d7261e\")", key, show(v))
		}
		ts.Accent = &s
		return nil
	}},
	"font_body":    familyKey(func(ts *Typesetting) **string { return &ts.FontBody }),
	"font_display": familyKey(func(ts *Typesetting) **string { return &ts.FontDisplay }),
	"font_mono":    familyKey(func(ts *Typesetting) **string { return &ts.FontMono }),
	"base_size": lengthKey(func(ts *Typesetting) **string { return &ts.BaseSize }, "1.0625rem",
		map[string][2]float64{"rem": {0.75, 1.5}, "px": {12, 24}}, "0.75rem to 1.5rem, or 12px to 24px"),
	"scale":   numberKey(func(ts *Typesetting) **float64 { return &ts.Scale }, 1.1, 1.618),
	"leading": numberKey(func(ts *Typesetting) **float64 { return &ts.Leading }, 1.1, 2.0),
	"measure": lengthKey(func(ts *Typesetting) **string { return &ts.Measure }, "34em",
		map[string][2]float64{"em": {20, 90}, "ch": {20, 90}}, "20 to 90 em or ch"),
	"align":            enumKey(func(ts *Typesetting) **string { return &ts.Align }, "left", "justify"),
	"grid":             enumKey(func(ts *Typesetting) **string { return &ts.Grid }, "asymmetric", "centered"),
	"margin_notes":     boolKey(func(ts *Typesetting) **bool { return &ts.MarginNotes }),
	"drop_caps":        enumKey(func(ts *Typesetting) **string { return &ts.DropCaps }, "none", "drop", "raised"),
	"lead_in":          enumKey(func(ts *Typesetting) **string { return &ts.LeadIn }, "none", "small_caps"),
	"numerals":         enumKey(func(ts *Typesetting) **string { return &ts.Numerals }, "roman", "arabic"),
	"rules":            enumKey(func(ts *Typesetting) **string { return &ts.Rules }, "none", "hairline", "heavy"),
	"fleurons":         boolKey(func(ts *Typesetting) **bool { return &ts.Fleurons }),
	"grayscale_images": boolKey(func(ts *Typesetting) **bool { return &ts.GrayscaleImages }),
}

// TypesettingKeys returns every valid `typesetting:` key, sorted.
func TypesettingKeys() []string {
	keys := make([]string, 0, len(typesettingKeys))
	for k := range typesettingKeys {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// ParseTypesetting validates a decoded `typesetting:` block and converts it
// into a Typesetting. Every problem is returned (sorted by key), each
// message naming the key and the accepted values — never just the first
// (contracts/config-and-tokens.md "Error format"). Unknown keys are errors
// so a misspelled setting is never silently ignored.
func ParseTypesetting(raw map[string]any) (Typesetting, []error) {
	var ts Typesetting
	keys := make([]string, 0, len(raw))
	for k := range raw {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var errs []error
	for _, key := range keys {
		spec, ok := typesettingKeys[key]
		if !ok {
			errs = append(errs, fmt.Errorf("typesetting.%s: unknown key (valid keys: %s)", key, strings.Join(TypesettingKeys(), ", ")))
			continue
		}
		if err := spec.parse(key, raw[key], &ts); err != nil {
			errs = append(errs, err)
		}
	}
	return ts, errs
}

// Resolve merges overrides onto the quarterly-default typesetting values.
// Only the quarterly preset has a typesetting schema; any other preset is
// an error.
func Resolve(preset Preset, overrides Typesetting) (Effective, error) {
	if preset.SiteType != "quarterly" {
		return Effective{}, fmt.Errorf("theme: preset %q has no typesetting schema (only quarterly presets do)", preset.ID)
	}
	eff := quarterlyDefaults
	var set []string
	str := func(key string, dst *string, v *string) {
		if v != nil {
			*dst = *v
			set = append(set, key)
		}
	}
	flt := func(key string, dst *float64, v *float64) {
		if v != nil {
			*dst = *v
			set = append(set, key)
		}
	}
	bln := func(key string, dst *bool, v *bool) {
		if v != nil {
			*dst = *v
			set = append(set, key)
		}
	}
	str("accent", &eff.Accent, overrides.Accent)
	str("font_body", &eff.FontBody, overrides.FontBody)
	str("font_display", &eff.FontDisplay, overrides.FontDisplay)
	str("font_mono", &eff.FontMono, overrides.FontMono)
	str("base_size", &eff.BaseSize, overrides.BaseSize)
	flt("scale", &eff.Scale, overrides.Scale)
	flt("leading", &eff.Leading, overrides.Leading)
	str("measure", &eff.Measure, overrides.Measure)
	str("align", &eff.Align, overrides.Align)
	str("grid", &eff.Grid, overrides.Grid)
	bln("margin_notes", &eff.MarginNotes, overrides.MarginNotes)
	str("drop_caps", &eff.DropCaps, overrides.DropCaps)
	str("lead_in", &eff.LeadIn, overrides.LeadIn)
	str("numerals", &eff.Numerals, overrides.Numerals)
	str("rules", &eff.Rules, overrides.Rules)
	bln("fleurons", &eff.Fleurons, overrides.Fleurons)
	bln("grayscale_images", &eff.GrayscaleImages, overrides.GrayscaleImages)
	if len(set) > 0 {
		sort.Strings(set)
		eff.Overridden = set
	}
	return eff, nil
}

func contains(values []string, s string) bool {
	for _, v := range values {
		if v == s {
			return true
		}
	}
	return false
}

// toFloat accepts the numeric shapes YAML decodes into interface{}.
func toFloat(v any) (float64, bool) {
	switch n := v.(type) {
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	case uint64:
		return float64(n), true
	case float64:
		if math.IsNaN(n) || math.IsInf(n, 0) {
			return 0, false
		}
		return n, true
	}
	return 0, false
}

// num formats a number the shortest way that round-trips (1.618, 2, 1.1).
func num(f float64) string {
	return strconv.FormatFloat(f, 'g', -1, 64)
}

// show renders an offending value for an error message: strings quoted,
// everything else as YAML would print it.
func show(v any) string {
	switch x := v.(type) {
	case string:
		return strconv.Quote(x)
	case float64:
		return num(x)
	case nil:
		return "null"
	}
	return fmt.Sprintf("%v", v)
}

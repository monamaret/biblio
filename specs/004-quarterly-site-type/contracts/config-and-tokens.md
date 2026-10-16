# Contract: settings and design values

## `content/site.yaml` → `typesetting:`

```yaml
site_type: quarterly
name: field-notes
description: ""
mcp_enabled: true
typesetting:            # optional; every key optional
  accent: "#1e4fd7"
  align: justify
  numerals: arabic
```

Keys, value types and defaults: the
[research doc's schema](../../research/004-quarterly-site-type.md#typesetting-schema-siteyaml).
Validation: [research.md R9](../research.md#r9-settings-validation-rules).

### Error format

One line per problem, all collected, then a non-zero exit:

```text
site.yaml: typesetting.scale: 2 is out of range (want a number from 1.1 to 1.618)
site.yaml: typesetting.align: "center" is not valid (want one of: left, justify)
site.yaml: typesetting.colour: unknown key (valid keys: accent, align, base_size, …)
site.yaml: typesetting is only supported for site_type quarterly (this site is "log")
```

### Warning format

```text
warning: typesetting.font_body: no file for "Söhne" in public/fonts/ — pages will use the fallback sans-serif
```

## `public/tokens.css`

- Generated; the header comment names the preset and lists the overridden
  keys. Deterministic (sorted names), so an unchanged config gives a
  byte-identical file.
- Only custom properties inside a single `:root { … }` block. No other
  selectors and no rules (R3).
- Quarterly emits the existing token names (`color-bg`, `color-fg`,
  `color-accent`, `color-muted`, `font-body`, `font-heading`) so shared
  conventions hold, plus:

| Property | From | Example (defaults) |
|---|---|---|
| `--color-rule`, `--color-plate-bg` | preset | `#111111`, `#ececec` |
| `--font-display`, `--font-mono` | `font_display`, `font_mono` | `"Archivo", "Helvetica Neue", Arial, sans-serif` |
| `--ts-base-size` | `base_size` | `1.0625rem` |
| `--ts-step-1` … `--ts-step-5` | `base_size` × `scale`ⁿ, computed in Go | `1.328rem` … |
| `--ts-leading` | `leading` | `1.5` |
| `--ts-measure` | `measure` | `34em` |
| `--ts-align`, `--ts-hyphens` | `align` | `left`, `manual` / `justify`, `auto` |
| `--ts-grid-columns` | `grid` | `minmax(0,2fr) minmax(0,7fr) minmax(0,3fr)` (asymmetric) |
| `--ts-rule-width` | `rules` | `0` / `1px` / `3px` |
| `--ts-dropcap-lines`, `--ts-dropcap-float` | `drop_caps` | `0`, `none` |
| `--ts-lead-in-variant` | `lead_in` | `all-small-caps` / `normal` |
| `--ts-numerals` | `numerals` | `upper-roman` / `decimal` |
| `--ts-fleuron` | `fleurons` | `"❦"` / `none` |

`font-heading` equals `font_display`. `--color-accent` equals `accent`.
`grayscale_images` and `margin_notes` produce no CSS: the first is used by
the image step, the second goes into the index.

## Template stylesheet `public/quarterly.css`

- Owned by the template, static, and linked from `root.tsx` after
  `/tokens.css`.
- Reads only `var(--…)` values. Outside `@font-face` it contains no color
  literal (`#…`, `rgb(`, `hsl(`, named colors) and no font family name.
  Enforced by `theme.test.tsx`.
- Contains `@font-face` rules for the bundled faces, the sidenote layout
  ([research.md R7](../research.md#r7-notes-in-html-without-client-js)), and
  `@media print` basics.

import type { Typesetting } from "./load";

const ROMAN: Array<[number, string]> = [
  [1000, "M"], [900, "CM"], [500, "D"], [400, "CD"], [100, "C"], [90, "XC"],
  [50, "L"], [40, "XL"], [10, "X"], [9, "IX"], [5, "V"], [4, "IV"], [1, "I"],
];

// toRoman renders a positive integer as upper-case roman numerals.
export function toRoman(n: number): string {
  let out = "";
  let rest = n;
  for (const [value, glyph] of ROMAN) {
    while (rest >= value) {
      out += glyph;
      rest -= value;
    }
  }
  return out;
}

// numeral renders an article's contents position in the site's numeral
// style — typesetting.numerals from the index, decided by biblio-cli
// (specs/004-quarterly-site-type research.md R3).
export function numeral(position: number, style: Typesetting["numerals"]): string {
  return style === "arabic" ? String(position) : toRoman(position);
}

// issueLine is the "No. 2 · Autumn 2026" label; the date stands in when an
// issue has no season (data-model.md "Issue" `season`).
export function issueLine(issue: { number: number; season?: string; date: string }): string {
  return `No. ${issue.number} · ${issue.season ?? issue.date}`;
}

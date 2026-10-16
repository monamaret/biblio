import contentIndex from "../../content-index.json";

// ContentEntry mirrors content-index.json's flat per-entry shape: slug/path
// plus every authored frontmatter field at the top level
// (internal/index/generator.go Entry.MarshalJSON).
export type ContentEntry = {
  slug: string;
  path: string;
  body: string;
  title?: string;
  [field: string]: unknown;
};

// listEntries/getEntryBySlug read content-index.json directly rather than
// re-parsing the authored Markdown+frontmatter on this side — biblio-cli's
// Go-side parser is the sole owner of that schema (Principle II: "Parsing,
// schema validation ... MUST NOT be duplicated between the Go CLI and the
// JS/TS site"). content-index.json already carries the parsed body and
// fields, so no separate `import.meta.glob`/frontmatter-parsing step is
// needed here.
const index = contentIndex as { site?: { name: string; description: string }; entries: ContentEntry[] };
const entries = index.entries;

export function listEntries(): ContentEntry[] {
  return entries;
}

export function getEntryBySlug(slug: string): ContentEntry | undefined {
  return entries.find((entry) => entry.slug === slug);
}

// siteMeta is site.yaml's name and description, copied into the index by
// biblio-cli. An index written before the `site` block existed has none.
export function siteMeta(): { name: string; description: string } {
  return index.site ?? { name: "", description: "" };
}

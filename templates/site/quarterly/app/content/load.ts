import contentIndex from "../../content-index.json";

// ContentEntry mirrors content-index.json's flat per-entry shape: slug/path
// plus every authored frontmatter field and the Go-derived fields at the
// top level (internal/index Entry.MarshalJSON).
export type ContentEntry = {
  slug: string;
  path: string;
  body: string;
  title?: string;
  kind?: "issue" | "article" | "colophon";
  [field: string]: unknown;
};

// Article is an entry of kind "article", with the fields biblio-cli derives
// from its issue folder (specs/004-quarterly-site-type data-model.md
// "Article").
export type Article = ContentEntry & {
  kind: "article";
  title: string;
  category: string;
  category_key: string;
  issue: string;
  issue_number: number;
  issue_date: string;
  position: number;
  author?: string;
  dek?: string;
  image?: string;
  image_alt?: string;
  image_caption?: string;
  plate?: string;
  tags?: string[];
};

export type IssueSummary = {
  issue: string;
  slug: string;
  number: number;
  title: string;
  date: string;
  season?: string;
  cover_plate?: string;
  cover_alt?: string;
  draft: boolean;
  article_count: number;
};

export type CategorySummary = { key: string; label: string; article_count: number };

export type Typesetting = { margin_notes: boolean; numerals: "roman" | "arabic" };

// QuarterlyIndex is the whole content-index.json a quarterly site ships
// (data-model.md "content-index.json shape").
export type QuarterlyIndex = {
  site?: { name: string; description: string };
  entries: ContentEntry[];
  issues: IssueSummary[];
  categories: CategorySummary[];
  current_issue: string | null;
  typesetting: Typesetting;
};

export type CategoryGroup = { category: { key: string; label: string }; articles: Article[] };

// Every accessor reads content-index.json as biblio-cli wrote it — issue
// membership, ordering, category keys and drafts are all decided Go-side,
// never re-derived from paths here (Principle II).
const index = contentIndex as unknown as QuarterlyIndex;

export function listEntries(): ContentEntry[] {
  return index.entries;
}

export function getEntryBySlug(slug: string): ContentEntry | undefined {
  return index.entries.find((entry) => entry.slug === slug);
}

export function listIssues(): IssueSummary[] {
  return index.issues;
}

// getIssue accepts an issue's folder slug ("02-autumn-2026") or its number
// as a string ("2").
export function getIssue(slugOrNumber: string): IssueSummary | undefined {
  return index.issues.find((i) => i.issue === slugOrNumber || String(i.number) === slugOrNumber);
}

export function currentIssue(): IssueSummary | undefined {
  return index.current_issue ? getIssue(index.current_issue) : undefined;
}

export function listCategories(): CategorySummary[] {
  return index.categories;
}

export function getCategory(key: string): CategorySummary | undefined {
  return index.categories.find((c) => c.key === key);
}

export function listArticles(): Article[] {
  return index.entries.filter((e): e is Article => e.kind === "article");
}

// articlesInIssue returns an issue's articles grouped by category, in the
// contents order biblio-cli already sorted them into.
export function articlesInIssue(issue: string): CategoryGroup[] {
  const groups: CategoryGroup[] = [];
  for (const article of listArticles().filter((a) => a.issue === issue)) {
    let group = groups.find((g) => g.category.key === article.category_key);
    if (!group) {
      const label = getCategory(article.category_key)?.label ?? article.category;
      group = { category: { key: article.category_key, label }, articles: [] };
      groups.push(group);
    }
    group.articles.push(article);
  }
  return groups;
}

export function articlesInCategory(key: string): Article[] {
  return listArticles().filter((a) => a.category_key === key);
}

export function getColophon(): ContentEntry | undefined {
  return index.entries.find((e) => e.kind === "colophon");
}

export function siteMeta(): { name: string; description: string } {
  return index.site ?? { name: "", description: "" };
}

export function typesetting(): Typesetting {
  return index.typesetting;
}

export function quarterlyIndex(): QuarterlyIndex {
  return index;
}

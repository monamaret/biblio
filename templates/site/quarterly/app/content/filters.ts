import { getIssue, listCategories, type Article } from "./load";

// ArticleQuery is list_articles' input: every field optional, combined with
// AND (specs/004-quarterly-site-type contracts/mcp-tools.md "list_articles").
export type ArticleQuery = {
  issue?: string;
  category?: string;
  author?: string;
  tag?: string;
  dateFrom?: string;
  dateTo?: string;
};

const DATE = /^\d{4}-\d{2}-\d{2}$/;

// filterArticles narrows articles (already in index order: newest issue
// first, then contents order) by the query. A category matches by its key
// or its label in any case — biblio-cli already folded every article's
// category into a key, so nothing is re-normalized here (Principle II).
// A malformed date is an error naming the field.
export function filterArticles(articles: Article[], query: ArticleQuery): { articles: Article[] } | { error: string } {
  for (const field of ["dateFrom", "dateTo"] as const) {
    const value = query[field];
    if (value !== undefined && !DATE.test(value)) {
      return { error: `Invalid ${field} "${value}": use YYYY-MM-DD.` };
    }
  }

  let categoryKey: string | undefined;
  if (query.category !== undefined) {
    const wanted = query.category.trim().toLowerCase();
    categoryKey = listCategories().find((c) => c.key === wanted || c.label.toLowerCase() === wanted)?.key ?? "\0";
  }
  const issue = query.issue !== undefined ? (getIssue(query.issue.trim())?.issue ?? "\0") : undefined;
  const author = query.author?.trim().toLowerCase();

  return {
    articles: articles.filter(
      (a) =>
        (issue === undefined || a.issue === issue) &&
        (categoryKey === undefined || a.category_key === categoryKey) &&
        (author === undefined || (a.author ?? "").toLowerCase() === author) &&
        (query.tag === undefined || (a.tags ?? []).includes(query.tag)) &&
        (query.dateFrom === undefined || a.issue_date >= query.dateFrom) &&
        (query.dateTo === undefined || a.issue_date <= query.dateTo),
    ),
  };
}

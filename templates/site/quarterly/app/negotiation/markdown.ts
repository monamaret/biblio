import {
  articlesInCategory,
  articlesInIssue,
  currentIssue,
  getCategory,
  getEntryBySlug,
  getIssue,
  listIssues,
  siteMeta,
  typesetting,
  type IssueSummary,
} from "../content/load";
import { issueLine, numeral } from "../content/numerals";

// renderMarkdown is tier-i: the Markdown view of any page, for clients that
// ask for text/markdown (and the base the ANSI tier renders from). List
// pages — the contents pages, the archive, category pages — are generated
// from the index; an article or the colophon is its title plus body
// (specs/004-quarterly-site-type contracts/site-surface.md "Content
// negotiation"). Returns undefined for a page that doesn't exist.
export function renderMarkdown(slug: string): string | undefined {
  if (slug === "") {
    const issue = currentIssue();
    if (!issue) {
      return `# ${siteMeta().name}\n\nNo issues have been published yet.\n`;
    }
    return contentsMarkdown(issue);
  }

  if (slug === "issues") {
    const lines = listIssues().map((issue) => {
      const count = `${issue.article_count} ${issue.article_count === 1 ? "article" : "articles"}`;
      return `- [${issueLine(issue)} — ${issue.title}](/issues/${issue.issue}) — ${count}`;
    });
    return `# Issues\n\n${lines.join("\n")}\n`;
  }

  const issueMatch = /^issues\/([^/]+)$/.exec(slug);
  if (issueMatch) {
    const issue = getIssue(issueMatch[1]);
    return issue && issue.issue === issueMatch[1] ? contentsMarkdown(issue) : undefined;
  }

  const categoryMatch = /^category\/([^/]+)$/.exec(slug);
  if (categoryMatch) {
    const category = getCategory(categoryMatch[1]);
    if (!category) return undefined;
    const lines = articlesInCategory(category.key).map((article) => {
      const issue = getIssue(article.issue);
      const author = article.author ? ` — ${article.author}` : "";
      return `- [${article.title}](/${article.slug})${author}${issue ? ` · ${issueLine(issue)}` : ""}`;
    });
    return `# ${category.label}\n\n${lines.join("\n")}\n`;
  }

  const entry = getEntryBySlug(slug);
  if (!entry || entry.kind === "issue") return undefined;
  const title = typeof entry.title === "string" ? entry.title : entry.slug;
  return `# ${title}\n\n${entry.body}`;
}

// contentsMarkdown is an issue's contents page as Markdown: the masthead,
// the issue line, then each category's articles with their numerals.
function contentsMarkdown(issue: IssueSummary): string {
  const style = typesetting().numerals;
  const parts = [`# ${siteMeta().name}`, `${issueLine(issue)} — ${issue.title}`, "## Contents"];
  const groups = articlesInIssue(issue.issue);
  if (groups.length === 0) {
    parts.push("No articles in this issue yet.");
  }
  for (const group of groups) {
    const lines = group.articles.map((article) => {
      const author = article.author ? ` — ${article.author}` : "";
      return `${numeral(article.position, style)}. [${article.title}](/${article.slug})${author}`;
    });
    parts.push(`### ${group.category.label}\n\n${lines.join("\n")}`);
  }
  parts.push("[Previous issues](/issues)");
  return parts.join("\n\n") + "\n";
}

import { z } from "zod";
import type { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import {
  articlesInIssue,
  getEntryBySlug,
  getIssue as findIssue,
  listArticles as allArticles,
  listIssues as allIssues,
  typesetting,
  type Article,
  type IssueSummary,
} from "../../content/load";
import { filterArticles, type ArticleQuery } from "../../content/filters";
import { issueLine, numeral } from "../../content/numerals";

// QuarterlyToolResult is the shared ToolResult shape plus the optional
// structured output and error flag the quarterly tools use.
export type QuarterlyToolResult = {
  content: Array<{ type: "text"; text: string }>;
  structuredContent?: Record<string, unknown>;
  isError?: boolean;
  [key: string]: unknown;
};

function textResult(text: string): QuarterlyToolResult {
  return { content: [{ type: "text", text }] };
}

type IndexNote = { id: string; label?: string; number?: number; text: string; ref_offset: number };

// Drafts are never visible to agents, even when `dev` keeps them in the
// index (contracts/mcp-tools.md).
function publishedIssues(): IssueSummary[] {
  return allIssues().filter((issue) => !issue.draft);
}

function publishedIssue(slugOrNumber: string): IssueSummary | undefined {
  const issue = findIssue(slugOrNumber.trim());
  return issue && !issue.draft ? issue : undefined;
}

function publishedArticles(): Article[] {
  return allArticles().filter((a) => !a.draft);
}

// listIssues: every published issue as a table, newest first.
export function listIssues(): QuarterlyToolResult {
  const issues = publishedIssues();
  if (issues.length === 0) {
    return textResult("No issues have been published yet.");
  }
  const rows = issues.map((i) => `| ${i.number} | ${i.title} | ${i.date} | ${i.article_count} | ${i.issue} |`);
  return textResult(["| No. | Title | Date | Articles | Slug |", "|---|---|---|---|---|", ...rows].join("\n"));
}

// getIssue: an issue's contents, grouped by category in contents order,
// numbered like the contents page.
export function getIssue(issue: string): QuarterlyToolResult {
  const found = publishedIssue(issue);
  if (!found) {
    return textResult(`No issue matches "${issue}". Use list_issues to see published issues.`);
  }
  const style = typesetting().numerals;
  const parts = [`# ${issueLine(found)} — ${found.title}`];
  const groups = articlesInIssue(found.issue);
  if (groups.length === 0) {
    parts.push("No articles in this issue yet.");
  }
  for (const group of groups) {
    const lines = group.articles.map((a) => {
      const author = a.author ? ` — ${a.author}` : "";
      return `${numeral(a.position, style)}. ${a.title}${author} (/${a.slug})`;
    });
    parts.push(`## ${group.category.label}\n\n${lines.join("\n")}`);
  }
  return textResult(parts.join("\n\n"));
}

// listArticles: article summaries matching every given filter.
export function listArticles(query: ArticleQuery): QuarterlyToolResult {
  const result = filterArticles(publishedArticles(), query);
  if ("error" in result) {
    return textResult(result.error);
  }
  if (result.articles.length === 0) {
    return textResult("No articles matched the given filter.");
  }
  const lines = result.articles.map((a) => {
    const author = a.author ? ` — ${a.author}` : "";
    const dek = a.dek ? ` — ${a.dek}` : "";
    return `- ${a.title}${author} · ${a.category} · No. ${a.issue_number}${dek} (/${a.slug})`;
  });
  return textResult(lines.join("\n"));
}

// normalizeSlug accepts "/issues/x/y", "issues/x/y" or "x/y".
function normalizeSlug(slug: string): string {
  const trimmed = slug.trim().replace(/^\/+/, "").replace(/\/+$/, "");
  return trimmed.startsWith("issues/") ? trimmed : `issues/${trimmed}`;
}

// getArticle: the full text (with margin notes as footnotes, so no note is
// lost in a terminal) plus the notes as structured data. MCP output is
// camelCase; the index is snake_case (contracts/mcp-tools.md "get_article").
export function getArticle(slug: string): QuarterlyToolResult {
  const entry = getEntryBySlug(normalizeSlug(slug));
  const article = entry?.kind === "article" && !entry.draft ? (entry as Article) : undefined;
  const issue = article ? publishedIssue(article.issue) : undefined;
  if (!article || !issue) {
    return { ...textResult(`No article matches "${slug}". Use list_articles to find one.`), isError: true };
  }

  const byline = [article.author, article.category, issueLine(issue)].filter(Boolean).join(" · ");
  const body = typeof article.markdown_body === "string" ? article.markdown_body : article.body;
  const footnotes = ((article.footnotes as IndexNote[]) ?? []).map((n) => ({
    id: n.id,
    label: n.label ?? "",
    number: n.number ?? 0,
    text: n.text,
    refOffset: n.ref_offset,
  }));
  const marginNotes = ((article.margin_notes as IndexNote[]) ?? []).map((n) => ({
    id: n.id,
    text: n.text,
    refOffset: n.ref_offset,
  }));

  return {
    ...textResult(`# ${article.title}\n\n${byline}\n\n${body}`),
    structuredContent: {
      slug: article.slug,
      title: article.title,
      ...(article.author ? { author: article.author } : {}),
      category: article.category,
      issue: { slug: issue.issue, number: issue.number, date: issue.date },
      footnotes,
      marginNotes,
    },
  };
}

const noteShape = { id: z.string(), text: z.string(), refOffset: z.number() };

// registerQuarterlyTools wires the four quarterly tools onto the server,
// next to the shared tools (contracts/mcp-tools.md).
export function registerQuarterlyTools(server: McpServer) {
  server.registerTool(
    "list_issues",
    { description: "List every published issue, newest first: number, title, date, article count, slug" },
    async () => listIssues(),
  );

  server.registerTool(
    "get_issue",
    {
      description: "Get an issue's contents, grouped by category, by its folder slug or its number",
      inputSchema: { issue: z.string() },
    },
    async ({ issue }: { issue: string }) => getIssue(issue),
  );

  server.registerTool(
    "list_articles",
    {
      description: "List articles, optionally filtered by issue, category, author, tag and issue-date range (all combined)",
      inputSchema: {
        issue: z.string().optional(),
        category: z.string().optional(),
        author: z.string().optional(),
        tag: z.string().optional(),
        dateFrom: z.string().optional(),
        dateTo: z.string().optional(),
      },
    },
    async (query: ArticleQuery) => listArticles(query),
  );

  server.registerTool(
    "get_article",
    {
      description: "Get an article's full text plus its footnotes and margin notes as structured data",
      inputSchema: { slug: z.string() },
      outputSchema: {
        slug: z.string(),
        title: z.string(),
        author: z.string().optional(),
        category: z.string(),
        issue: z.object({ slug: z.string(), number: z.number(), date: z.string() }),
        footnotes: z.array(z.object({ ...noteShape, label: z.string(), number: z.number() })),
        marginNotes: z.array(z.object(noteShape)),
      },
    },
    async ({ slug }: { slug: string }) => getArticle(slug),
  );
}

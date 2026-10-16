import fs from "node:fs";
import path from "node:path";
import type { Plugin } from "vite";
import type { QuarterlyIndex } from "../content/load";
import { issueLine } from "../content/numerals";

// buildLlmsTxt produces the minimal, build-time llms.txt for a quarterly
// site: the site name and description, then the current issue, every
// issue, the current issue's articles, and the colophon — links only,
// never bodies, never drafts (specs/004-quarterly-site-type
// contracts/site-surface.md "llms.txt").
export function buildLlmsTxt(index: QuarterlyIndex): string {
  const site = index.site ?? { name: "", description: "" };
  const issues = index.issues.filter((issue) => !issue.draft);
  const current = issues.find((issue) => issue.issue === index.current_issue);

  const lines = [`# ${site.name}`, "", site.description, ""];

  if (current) {
    lines.push("## Current issue", "", `- ${issueLine(current)}: /issues/${current.issue}`, "");
  }

  if (issues.length > 0) {
    lines.push("## Issues", "");
    for (const issue of issues) {
      lines.push(`- ${issueLine(issue)}: /issues/${issue.issue}`);
    }
    lines.push("");
  }

  if (current) {
    const articles = index.entries.filter((e) => e.kind === "article" && e.issue === current.issue && !e.draft);
    if (articles.length > 0) {
      lines.push("## Articles (current issue)", "");
      for (const article of articles) {
        const author = typeof article.author === "string" ? ` — ${article.author}` : "";
        lines.push(`- ${article.title}${author}: /${article.slug}`);
      }
      lines.push("");
    }
  }

  if (index.entries.some((e) => e.kind === "colophon")) {
    lines.push("## About", "", "- Colophon: /colophon", "");
  }

  return lines.join("\n");
}

// llmsTxtPlugin is a Vite plugin writing llms.txt into the client build
// output during `biblio-cli build`, derived from the project's
// content-index.json — which biblio-cli has just regenerated, site name and
// description included, so nothing here parses site.yaml.
export function llmsTxtPlugin(): Plugin {
  return {
    name: "biblio-llms-txt",
    apply: "build",
    closeBundle() {
      const root = process.cwd();
      const contentIndexPath = path.join(root, "content-index.json");
      const outDir = path.join(root, "dist", "client");
      if (!fs.existsSync(contentIndexPath) || !fs.existsSync(outDir)) {
        return;
      }
      const index = JSON.parse(fs.readFileSync(contentIndexPath, "utf-8")) as QuarterlyIndex;
      fs.writeFileSync(path.join(outDir, "llms.txt"), buildLlmsTxt(index));
    },
  };
}

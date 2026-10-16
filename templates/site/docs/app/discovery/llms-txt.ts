import fs from "node:fs";
import path from "node:path";
import type { Plugin } from "vite";
import type { ContentEntry } from "../content/load";

// buildLlmsTxt produces the minimal, build-time-generated llms.txt content:
// site_name/description plus a flat list of key page links, derived
// directly from content-index.json — not a fuller `llms-full.txt`-style
// content dump (research.md "llms.txt investment level"; FR-005).
export function buildLlmsTxt(siteName: string, description: string, entries: ContentEntry[]): string {
  // An empty description (what `biblio-cli new` writes) gets no paragraph
  // at all rather than a blank one (#217).
  const lines = [`# ${siteName}`, ""];
  if (description.trim() !== "") {
    lines.push(description, "");
  }
  lines.push("## Pages", "");
  for (const entry of entries) {
    const title = typeof entry.title === "string" ? entry.title : entry.slug;
    lines.push(`- [${title}](/${entry.slug})`);
  }
  return lines.join("\n") + "\n";
}

// readSiteMeta extracts just `name`/`description` from content/site.yaml —
// a fixed two-field read, not a general YAML parser, since adding a YAML
// dependency for this single derived discovery artifact would be exactly
// the kind of speculative dependency Principle I rules out (the canonical
// site.yaml schema parsing already lives Go-side, per Principle II).
function readSiteMeta(siteYamlPath: string): { name: string; description: string } {
  const raw = fs.readFileSync(siteYamlPath, "utf-8");
  const name = /^name:\s*(.*)$/m.exec(raw)?.[1]?.trim() ?? "";
  const description = /^description:\s*"?(.*?)"?$/m.exec(raw)?.[1]?.trim() ?? "";
  return { name, description };
}

// llmsTxtPlugin is a Vite plugin writing llms.txt into the client build
// output during `biblio-cli build` (a thin wrapper over `vite build`),
// derived from the project's content-index.json and content/site.yaml.
export function llmsTxtPlugin(): Plugin {
  return {
    name: "biblio-llms-txt",
    apply: "build",
    closeBundle() {
      const root = process.cwd();
      const contentIndexPath = path.join(root, "content-index.json");
      const siteYamlPath = path.join(root, "content", "site.yaml");
      if (!fs.existsSync(contentIndexPath) || !fs.existsSync(siteYamlPath)) {
        return;
      }
      const { entries } = JSON.parse(fs.readFileSync(contentIndexPath, "utf-8")) as {
        entries: ContentEntry[];
      };
      const { name, description } = readSiteMeta(siteYamlPath);
      const outDir = path.join(root, "dist", "client");
      if (!fs.existsSync(outDir)) {
        return;
      }
      fs.writeFileSync(path.join(outDir, "llms.txt"), buildLlmsTxt(name, description, entries));
    },
  };
}

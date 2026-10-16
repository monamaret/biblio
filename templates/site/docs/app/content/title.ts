import { siteMeta } from "./load";

// pageTitle is a page's <title>: "On Margins · 502.io" for a page, the bare
// site name for the home page, and just the page when site.yaml has no
// name (#215).
export function pageTitle(page?: string): string {
  const site = siteMeta().name.trim();
  if (!page) {
    return site;
  }
  return site ? `${page} · ${site}` : page;
}

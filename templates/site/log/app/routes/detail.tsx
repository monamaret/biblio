import { Link, useParams, type MetaFunction } from "react-router";
import { getEntryBySlug } from "../content/load";
import { ContentBody } from "../content/render";
import { pageTitle } from "../content/title";

// The page title is the entry's title, then the site name (#215).
export const meta: MetaFunction = ({ params }) => {
  const entry = getEntryBySlug((params["*"] ?? "").replace(/\/+$/, ""));
  const title = entry ? (typeof entry.title === "string" ? entry.title : entry.slug) : "Not found";
  return [{ title: pageTitle(title) }];
};

// Detail is the per-entry view, displaying date/tags/project (project
// omitted when absent — it's optional, data-model.md "Project-Log Entry")
// plus the rendered Markdown body (FR-002; spec Scenario 1).
export default function Detail() {
  const params = useParams();
  // Cloudflare's static-asset serving redirects "/foo" to "/foo/" for
  // directory-style prerendered output, so a hard refresh's splat param
  // can carry a trailing slash that no content-index.json slug has.
  const slug = (params["*"] ?? "").replace(/\/+$/, "");
  const entry = getEntryBySlug(slug);

  if (!entry) {
    return (
      <main>
        <h1>Not found</h1>
        <p>
          No entry matches <code>{slug}</code>.
        </p>
        <Link to="/">Back to the log</Link>
      </main>
    );
  }

  const tags = Array.isArray(entry.tags) ? (entry.tags as unknown[]) : [];

  return (
    <main>
      <h1>{typeof entry.title === "string" ? entry.title : entry.slug}</h1>
      <p>
        {typeof entry.date === "string" && <span>{entry.date}</span>}
        {typeof entry.project === "string" && <span> — {entry.project}</span>}
        {tags.length > 0 && <span> — {tags.join(", ")}</span>}
      </p>
      <ContentBody body={entry.body} />
      <p>
        <Link to="/">Back to the log</Link>
      </p>
    </main>
  );
}

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

// Detail is the per-page view for any documentation page — one route per
// content-index.json entry, matching FR-002's completeness invariant
// (spec Scenario 1.1). Displays the page's section context alongside its
// rendered Markdown body.
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
          No page matches <code>{slug}</code>.
        </p>
        <Link to="/">Back home</Link>
      </main>
    );
  }

  const section = typeof entry.section === "string" ? entry.section : undefined;

  return (
    <main>
      {section && <p>{section}</p>}
      <h1>{typeof entry.title === "string" ? entry.title : entry.slug}</h1>
      <ContentBody body={entry.body} />
      <p>
        <Link to="/">Back home</Link>
      </p>
    </main>
  );
}

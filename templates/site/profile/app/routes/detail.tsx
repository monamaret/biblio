import { Link, useParams } from "react-router";
import { getEntryBySlug } from "../content/load";
import { ContentBody } from "../content/render";

// Detail is the per-entry view for any content-driven page (profile,
// project, or future experience/skills items) — one route per
// content-index.json entry, matching FR-001's completeness invariant.
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
          No content entry matches <code>{slug}</code>.
        </p>
        <Link to="/">Back home</Link>
      </main>
    );
  }

  return (
    <main>
      <h1>{typeof entry.title === "string" ? entry.title : entry.slug}</h1>
      <ContentBody body={entry.body} />
      <p>
        <Link to="/">Back home</Link>
      </p>
    </main>
  );
}

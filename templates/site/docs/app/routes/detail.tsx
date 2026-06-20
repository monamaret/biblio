import { Link, useParams } from "react-router";
import { getEntryBySlug } from "../content/load";
import { ContentBody } from "../content/render";

// Detail is the per-page view for any documentation page — one route per
// content-index.json entry, matching FR-002's completeness invariant
// (spec Scenario 1.1). Displays the page's section context alongside its
// rendered Markdown body.
export default function Detail() {
  const params = useParams();
  const slug = params["*"] ?? "";
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

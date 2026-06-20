import { Link } from "react-router";
import { listEntries } from "../content/load";
import { toc } from "../content/toc";

// Home is the TOC/navigation view (FR-002; spec Scenario 1.1) — every
// authored page, grouped under its section heading and ordered within it.
// Deliberately has no tag/project filter controls anywhere (Scenario 1.2;
// contracts/docs-site-and-review.md "Filter UI") — that absence is this
// site type's confirmed design, not an unfinished feature.
export default function Home() {
  const sections = toc(listEntries());

  return (
    <main>
      <h1>Welcome to your biblio docs site</h1>
      {sections.map((section) => (
        <section key={section.section}>
          <h2>{section.section}</h2>
          <nav>
            <ul>
              {section.pages.map((page) => (
                <li key={page.slug}>
                  <Link to={`/${page.slug}`}>
                    {typeof page.title === "string" ? page.title : page.slug}
                  </Link>
                </li>
              ))}
            </ul>
          </nav>
        </section>
      ))}
    </main>
  );
}

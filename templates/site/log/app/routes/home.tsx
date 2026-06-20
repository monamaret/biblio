import { Link } from "react-router";
import { listEntries } from "../content/load";

// Home is the chronological entry-list view (FR-002; spec Scenario 1) —
// every authored entry, newest first, with its date and tags (and project,
// when present — `project` is optional, data-model.md "Project-Log Entry")
// shown alongside the title link to its detail route. Filter controls are
// added in Phase 4 (T026).
export default function Home() {
  const entries = [...listEntries()].sort((a, b) => {
    const aDate = typeof a.date === "string" ? a.date : "";
    const bDate = typeof b.date === "string" ? b.date : "";
    return bDate.localeCompare(aDate);
  });

  return (
    <main>
      <h1>Project Log</h1>
      <ul>
        {entries.map((entry) => {
          const tags = Array.isArray(entry.tags) ? (entry.tags as unknown[]) : [];
          return (
            <li key={entry.slug}>
              <Link to={`/${entry.slug}`}>
                {typeof entry.title === "string" ? entry.title : entry.slug}
              </Link>
              {typeof entry.date === "string" && <span> — {entry.date}</span>}
              {typeof entry.project === "string" && <span> — {entry.project}</span>}
              {tags.length > 0 && <span> — {tags.join(", ")}</span>}
            </li>
          );
        })}
      </ul>
    </main>
  );
}

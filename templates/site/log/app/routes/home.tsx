import { useState } from "react";
import { Link } from "react-router";
import { listEntries } from "../content/load";
import { filterEntries, type FilterQuery } from "../content/filters";

// Home is the chronological, filterable entry-list view (FR-002, FR-003;
// spec Scenario 1 & 2) — every authored entry, newest first, with its
// date and tags (and project, when present — `project` is optional,
// data-model.md "Project-Log Entry") shown alongside the title link to its
// detail route. The date-range/tag/project controls below feed
// `filterEntries()` (the canonical, single filter-computation module also
// used by the log-specific MCP tools); options are derived directly from
// whatever values exist on the entries, so a brand-new tag/project is
// filterable immediately after the next build with no registration step
// (spec Edge Cases).
export default function Home() {
  const allEntries = listEntries();
  const [query, setQuery] = useState<FilterQuery>({});

  const tagOptions = Array.from(
    new Set(allEntries.flatMap((entry) => (Array.isArray(entry.tags) ? (entry.tags as string[]) : []))),
  ).sort();
  const projectOptions = Array.from(
    new Set(
      allEntries
        .map((entry) => (typeof entry.project === "string" ? entry.project : undefined))
        .filter((project): project is string => Boolean(project)),
    ),
  ).sort();

  const entries = [...filterEntries(allEntries, query)].sort((a, b) => {
    const aDate = typeof a.date === "string" ? a.date : "";
    const bDate = typeof b.date === "string" ? b.date : "";
    return bDate.localeCompare(aDate);
  });

  return (
    <main>
      <h1>Project Log</h1>

      <form aria-label="Filter entries">
        <label>
          From
          <input
            type="date"
            aria-label="Date from"
            value={query.dateFrom ?? ""}
            onChange={(e) => setQuery({ ...query, dateFrom: e.target.value || undefined })}
          />
        </label>
        <label>
          To
          <input
            type="date"
            aria-label="Date to"
            value={query.dateTo ?? ""}
            onChange={(e) => setQuery({ ...query, dateTo: e.target.value || undefined })}
          />
        </label>
        <label>
          Tag
          <select
            aria-label="Tag"
            value={query.tag ?? ""}
            onChange={(e) => setQuery({ ...query, tag: e.target.value || undefined })}
          >
            <option value="">All tags</option>
            {tagOptions.map((tag) => (
              <option key={tag} value={tag}>
                {tag}
              </option>
            ))}
          </select>
        </label>
        <label>
          Project
          <select
            aria-label="Project"
            value={query.project ?? ""}
            onChange={(e) => setQuery({ ...query, project: e.target.value || undefined })}
          >
            <option value="">All projects</option>
            {projectOptions.map((project) => (
              <option key={project} value={project}>
                {project}
              </option>
            ))}
          </select>
        </label>
      </form>

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

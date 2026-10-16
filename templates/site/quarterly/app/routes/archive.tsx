import { Link, type MetaFunction } from "react-router";
import { listIssues } from "../content/load";
import { Masthead } from "../content/masthead";
import { issueLine } from "../content/numerals";
import { pageTitle } from "../content/title";

// The archive is titled "Issues · <site name>" (#215).
export const meta: MetaFunction = () => [{ title: pageTitle("Issues") }];

// Archive lists every published issue, newest first (contracts/site-surface.md
// "Routes" `/issues`).
export default function Archive() {
  const issues = listIssues();
  return (
    <main className="archive-page">
      <Masthead />
      <h1>Issues</h1>
      {issues.length === 0 ? (
        <p className="empty">No issues have been published yet.</p>
      ) : (
        <ol className="archive">
          {issues.map((issue) => (
            <li className="archive-item" key={issue.issue}>
              <Link to={`/issues/${issue.issue}`}>
                {issueLine(issue)} — {issue.title}
              </Link>{" "}
              <span className="article-count">
                {issue.article_count} {issue.article_count === 1 ? "article" : "articles"}
              </span>
            </li>
          ))}
        </ol>
      )}
    </main>
  );
}

import { Link, useParams, type MetaFunction } from "react-router";
import { articlesInIssue, currentIssue, getIssue, typesetting } from "../content/load";
import { Masthead, NotFound } from "../content/masthead";
import { issueLine, numeral } from "../content/numerals";
import { pageTitle } from "../content/title";

// "/" is titled with the site name alone; an issue page with the issue's title (#215).
export const meta: MetaFunction = ({ params }) => {
  if (!params.issue) {
    return [{ title: pageTitle() }];
  }
  return [{ title: pageTitle(getIssue(params.issue)?.title ?? "Not found") }];
};

// Contents is the issue contents page — the current issue at "/", any
// published issue at "/issues/:issue": masthead, "No. N · season", cover
// plate, then articles grouped by category with numerals, dotted leaders
// and authors (specs/004-quarterly-site-type FR-013, FR-014;
// contracts/site-surface.md "Contents page").
export default function Contents() {
  const params = useParams();
  const issue = params.issue ? getIssue(params.issue) : currentIssue();

  if (params.issue && !issue) {
    return <NotFound what={`/issues/${params.issue}`} />;
  }

  if (!issue) {
    return (
      <main>
        <Masthead />
        <p className="empty">No issues have been published yet.</p>
      </main>
    );
  }

  const groups = articlesInIssue(issue.issue);
  const style = typesetting().numerals;

  return (
    <main className="contents-page">
      <Masthead />
      <p className="issue-line">{issueLine(issue)}</p>
      <h1>{issue.title}</h1>
      {issue.cover_plate && (
        <figure className="plate cover-plate">
          <img src={issue.cover_plate} alt={issue.cover_alt ?? ""} />
        </figure>
      )}
      <h2>Contents</h2>
      {groups.length === 0 ? (
        <p className="empty">No articles in this issue yet.</p>
      ) : (
        groups.map((group) => (
          <section className="contents-category" key={group.category.key}>
            <h3>{group.category.label}</h3>
            <ol className="contents-list">
              {group.articles.map((article) => (
                <li className="contents-item" key={article.slug}>
                  <span className="numeral">{numeral(article.position, style)}.</span>{" "}
                  <Link to={`/${article.slug}`}>{article.title}</Link>
                  <span className="leader" aria-hidden="true" />
                  {article.author && <span className="author">{article.author}</span>}
                </li>
              ))}
            </ol>
          </section>
        ))
      )}
      <p className="previous">
        <Link to="/issues">Previous issues →</Link>
      </p>
    </main>
  );
}

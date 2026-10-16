import { Link, useParams } from "react-router";
import { articlesInCategory, getCategory, getIssue } from "../content/load";
import { Masthead, NotFound } from "../content/masthead";
import { issueLine } from "../content/numerals";

// Category lists one category's articles across every issue, newest issue
// first, each showing its issue (specs/004-quarterly-site-type FR-011,
// FR-013).
export default function Category() {
  const params = useParams();
  const key = params.category ?? "";
  const category = getCategory(key);
  if (!category) {
    return <NotFound what={`/category/${key}`} />;
  }

  return (
    <main className="category-page">
      <Masthead />
      <h1>{category.label}</h1>
      <ol className="category-list">
        {articlesInCategory(key).map((article) => {
          const issue = getIssue(article.issue);
          return (
            <li className="category-item" key={article.slug}>
              <Link to={`/${article.slug}`}>{article.title}</Link>
              {article.author && <span className="author"> — {article.author}</span>}
              {issue && (
                <>
                  {" · "}
                  <Link to={`/issues/${issue.issue}`}>{issueLine(issue)}</Link>
                </>
              )}
            </li>
          );
        })}
      </ol>
    </main>
  );
}

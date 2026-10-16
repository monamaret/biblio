import { Link, useParams } from "react-router";
import { getEntryBySlug, getIssue, type Article } from "../content/load";
import { ContentBody } from "../content/render";
import { Masthead, NotFound } from "../content/masthead";
import { issueLine } from "../content/numerals";

// Detail renders an article or the colophon (contracts/site-surface.md
// "Article page"): category label, title, author, dek, the lead plate as
// "Fig. 1", then the body. Issue entries are served by the contents route,
// never here.
export default function Detail() {
  const params = useParams();
  const slug = params["*"] ?? "";
  const entry = getEntryBySlug(slug);

  if (!entry || entry.kind === "issue") {
    return <NotFound what={`/${slug}`} />;
  }

  if (entry.kind === "colophon") {
    return (
      <main className="colophon-page">
        <Masthead />
        <h1>{entry.title}</h1>
        <div className="article-body">
          <ContentBody body={entry.body} />
        </div>
      </main>
    );
  }

  const article = entry as Article;
  const issue = getIssue(article.issue);

  return (
    <main className="article-page">
      <Masthead />
      <article>
        <p className="category-label">
          <Link to={`/category/${article.category_key}`}>{article.category}</Link>
        </p>
        <h1>{article.title}</h1>
        {article.author && <p className="byline">{article.author}</p>}
        {article.dek && <p className="dek">{article.dek}</p>}
        {article.plate && (
          <figure className="plate lead-plate">
            <img src={article.plate} alt={article.image_alt ?? ""} />
            {article.image_caption && <figcaption>Fig. 1 — {article.image_caption}</figcaption>}
          </figure>
        )}
        <div className="article-body">
          <ContentBody body={article.body} />
        </div>
      </article>
      {issue && (
        <p className="issue-link">
          <Link to={`/issues/${issue.issue}`}>{issueLine(issue)}</Link>
        </p>
      )}
    </main>
  );
}

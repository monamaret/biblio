import { Link } from "react-router";
import { getColophon, siteMeta } from "./load";

// Masthead is shared by every page: the site name, the archive, and the
// colophon only when one exists (contracts/site-surface.md "Routes").
export function Masthead() {
  const colophon = getColophon();
  return (
    <header className="masthead">
      <p className="masthead-title">
        <Link to="/">{siteMeta().name}</Link>
      </p>
      <nav className="masthead-nav">
        <Link to="/issues">Issues</Link>
        {colophon && <Link to="/colophon">Colophon</Link>}
      </nav>
    </header>
  );
}

// NotFound is the in-page not-found view every route renders for an
// unknown issue, category or slug, instead of throwing.
export function NotFound({ what }: { what: string }) {
  return (
    <main>
      <Masthead />
      <h1>Not found</h1>
      <p>
        Nothing matches <code>{what}</code>.
      </p>
      <p>
        <Link to="/">Back to the current issue</Link>
      </p>
    </main>
  );
}

import { Link } from "react-router";
import { listEntries } from "../content/load";

// Home is the list/navigation view (FR-001, FR-003; spec Scenario 1 & 3) —
// every authored entry, linked to its detail route.
export default function Home() {
  const entries = listEntries();

  return (
    <main>
      <h1>Welcome to your biblio site</h1>
      <nav>
        <ul>
          {entries.map((entry) => (
            <li key={entry.slug}>
              <Link to={`/${entry.slug}`}>
                {typeof entry.title === "string" ? entry.title : entry.slug}
              </Link>
            </li>
          ))}
        </ul>
      </nav>
    </main>
  );
}

import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import contentIndex from "../content-index.json";

type Entry = { slug: string; title?: string };

function App() {
  const entries = (contentIndex as { entries: Entry[] }).entries;
  return (
    <main>
      <h1>Welcome to your biblio site</h1>
      <ul>
        {entries.map((entry) => (
          <li key={entry.slug}>{entry.title ?? entry.slug}</li>
        ))}
      </ul>
    </main>
  );
}

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <App />
  </StrictMode>,
);

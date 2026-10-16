import type { ReactNode } from "react";
import { Links, Meta, Outlet, Scripts, ScrollRestoration } from "react-router";

// Layout applies the unchanged v0.01-generated tokens.css theme to every
// route — the one styling path every view renders through (data-model.md
// "Deployment" `theme`).
export function Layout({ children }: { children: ReactNode }) {
  return (
    <html lang="en">
      <head>
        <meta charSet="UTF-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1.0" />
        <link rel="stylesheet" href="/tokens.css" />
        <link rel="stylesheet" href="/quarterly.css" />
        {/* Every page is discoverable in its Markdown-friendly form at the
            same address (content-negotiated via Accept: text/markdown —
            app/negotiation/accept.ts) — this hint surfaces that to
            agents/crawlers without needing a separate sibling URL (FR-005). */}
        <link rel="alternate" type="text/markdown" href="." />
        <link rel="help" type="text/plain" href="/llms.txt" title="Site summary for agents/crawlers" />
        <Meta />
        <Links />
        <title>biblio quarterly site</title>
      </head>
      <body>
        {children}
        <ScrollRestoration />
        <Scripts />
      </body>
    </html>
  );
}

export default function App() {
  return <Outlet />;
}

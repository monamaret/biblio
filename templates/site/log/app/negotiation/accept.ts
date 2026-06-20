import { renderMarkdown } from "./markdown";
import { renderToAnsi } from "./ansi";

export type Variant = "html" | "markdown" | "ansi";

// TERMINAL_UA_PATTERN names known terminal/CLI clients that can render
// ANSI escape codes — anything else preferring Markdown gets the plain
// tier-i passthrough instead (research.md "Deterministic Accept/User-Agent
// tie-breaking rule").
const TERMINAL_UA_PATTERN = /\b(curl|wget|httpie|lynx)\b/i;

// qValueFor returns the parsed `q` weight (default 1) for the given media
// type if it's present anywhere in the Accept header, or undefined if it
// isn't mentioned at all. A minimal parser — good enough for the small,
// fixed set of media types this negotiation layer cares about.
function qValueFor(accept: string | null | undefined, mediaType: string): number | undefined {
  if (!accept) return undefined;
  for (const part of accept.split(",")) {
    const [rawType, ...params] = part.trim().split(";");
    if (rawType.trim().toLowerCase() !== mediaType) continue;
    for (const param of params) {
      const [key, value] = param.trim().split("=");
      if (key.trim().toLowerCase() === "q") {
        const q = Number.parseFloat(value);
        return Number.isNaN(q) ? 1 : q;
      }
    }
    return 1;
  }
  return undefined;
}

// negotiateVariant resolves the documented, deterministic precedence rule
// (research.md "Deterministic Accept/User-Agent tie-breaking rule"):
// `text/markdown` wins over `text/html` whenever it's mentioned at all
// with an equal-or-higher `q` (including an exact tie); otherwise HTML —
// the unchanged v0.02 browser-facing default. Within the markdown family,
// a recognized terminal User-Agent selects ANSI (tier ii); anything else
// gets plain Markdown (tier i).
export function negotiateVariant(accept: string | null | undefined, userAgent: string | null | undefined): Variant {
  const markdownQ = qValueFor(accept, "text/markdown");
  if (markdownQ === undefined) {
    return "html";
  }
  const htmlQ = qValueFor(accept, "text/html") ?? qValueFor(accept, "*/*") ?? 0;
  if (markdownQ < htmlQ) {
    return "html";
  }
  return TERMINAL_UA_PATTERN.test(userAgent ?? "") ? "ansi" : "markdown";
}

const VARY_HEADER = "Accept, User-Agent";

// negotiatePageResponse decides the response for a page request, given its
// content-index slug. Returns `undefined` when the negotiated variant is
// `html` — the caller (the Worker entry script) should fall through to
// Workers Static Assets/the prerendered HTML for that case, unchanged from
// v0.02. Every response this function actually returns carries
// `Vary: Accept, User-Agent` so Cloudflare's edge cache serves the correct
// variant to the correct client (research.md "Negotiation and caching
// correctness"). A slug with no matching content-index entry gets a clear,
// structured 404 — never a silent empty/200 response (spec Edge Cases).
export async function negotiatePageResponse(
  slug: string,
  accept: string | null | undefined,
  userAgent: string | null | undefined,
): Promise<Response | undefined> {
  const variant = negotiateVariant(accept, userAgent);
  if (variant === "html") {
    return undefined;
  }

  const markdown = renderMarkdown(slug);
  if (markdown === undefined) {
    return new Response(`Not found: no page matches "${slug}"`, {
      status: 404,
      headers: { "Content-Type": "text/plain; charset=utf-8", Vary: VARY_HEADER },
    });
  }

  if (variant === "ansi") {
    const ansi = await renderToAnsi(markdown);
    return new Response(ansi, {
      headers: { "Content-Type": "text/plain; charset=utf-8", Vary: VARY_HEADER },
    });
  }

  return new Response(markdown, {
    headers: { "Content-Type": "text/markdown; charset=utf-8", Vary: VARY_HEADER },
  });
}

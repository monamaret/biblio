# Architecture: `biblio` platform — reference diagrams

Reference diagrams for the `biblio` architecture selected in [001-static-site-platform-with-mcp.md](001-static-site-platform-with-mcp.md) — see that document for the research and rationale behind these decisions, and the [Decision Log](001-static-site-platform-with-mcp.md#decision-log) for exact wording of each confirmed decision.

## 1. High-level architecture (selected)

A single Cloudflare Worker is the **one public origin** for a `biblio` site: it serves the prerendered React + Vite static build *and* implements the `/mcp` endpoint in-Worker, on one domain/host with no cross-origin boundary.

```mermaid
flowchart TB
    visitor["Browser visitor"]
    mcpclient["MCP client / AI assistant"]
    terminal["Terminal client<br/>curl / lynx / coding agent"]

    subgraph CF["Cloudflare (one zone, custom domain)"]
        direction TB
        worker["Worker: example.com<br/>fetch handler — routes by path + Accept"]
        assets[("Workers Static Assets<br/>prerendered site bundle")]
        mcp["In-Worker MCP server<br/>agents / workers-mcp SDK"]
        do[("Durable Objects<br/>MCP session state")]
        index[["content-index.json<br/>inlined in Worker bundle"]]
        d1[("D1 — escalation only")]
        vec[("Vectorize — escalation only")]
    end

    visitor -->|"GET /*"| worker
    terminal -->|"GET /path, Accept: text/markdown<br/>or /path.md"| worker
    mcpclient -->|"GET/POST /mcp"| worker

    worker -->|"/* (HTML request)"| assets
    worker -->|"/* (markdown-negotiated)"| index
    worker -->|"/mcp*"| mcp
    mcp --> do
    mcp --> index
    mcp -.->|"escalation path"| d1
    mcp -.->|"escalation path"| vec
```

**Why this shape:** the site and `/mcp` share one domain/host (no CORS to manage), one deploy, one set of bindings — the simplest topology to scaffold and replicate across many sites. See §5 below for the alternative patterns that were evaluated and dropped.

## 2. Build & deployment pipeline

```mermaid
flowchart LR
    content[/"Content directory<br/>content/&lt;site-type&gt;/*.md"/]
    cli["biblio-cli<br/>Go-side generator"]
    index[["content-index.json"]]
    tokens[["public/tokens.css<br/>(quarterly)"]]
    plates[["public/plates/<br/>(quarterly)"]]
    rr["react-router build<br/>(prerender via react-router.config.ts)"]
    dist[/"Prerendered static bundle<br/>HTML + CSS + JS"/]
    wrangler["Wrangler build & deploy"]
    worker(("Cloudflare Worker<br/>+ bindings"))
    gh["GitHub Actions<br/>(on push)"]

    gh --> cli
    content --> cli --> index
    cli -.-> tokens
    cli -.-> plates
    index --> rr
    tokens -.-> rr
    plates -.-> rr
    content --> rr
    rr --> dist
    index --> wrangler
    dist --> wrangler
    wrangler --> worker
```

Three build steps run in order, normally triggered by a push via GitHub Actions:

0. **Pre-build step inside `biblio-cli build`/`dev`** *(planned, [§11](roadmap.md#11-quarterly-site-type), [004](004-quarterly-site-type.md#build-pipeline-changes))* — for every site type, validate `content/site.yaml` and run step 1 automatically, so `index` no longer has to be run by hand. For `quarterly` sites only, it also regenerates `public/tokens.css` from the preset plus the validated `typesetting:` overrides, and writes grayscale, width-capped copies of content images to `public/plates/` (dashed edges above). Order: validate → tokens → plates → index. Draft issues are left out of the index on `build` and kept on `dev`.
1. **`biblio-cli` generates `content-index.json`** — a Go-side generator parses every Markdown file's frontmatter + body (mirroring what `gray-matter` would do in JS) into a single canonical artifact, the source of truth for both the site build and the MCP server.
2. **`react-router build`** — the framework's native `prerender` (configured in `react-router.config.ts`, driven by an async function that enumerates content-driven routes from the content index) renders every route to static HTML/CSS/JS.
3. **Wrangler build & deploy** — the Worker `import`s `content-index.json` directly (inlined into its V8 bundle for in-memory querying, no extra fetch at runtime) and binds the prerendered output as static assets, then publishes both to Cloudflare.

## 3. Request routing call flow (in-Worker `fetch` handler)

```mermaid
sequenceDiagram
    participant C as Client
    participant W as Worker (fetch handler)
    participant A as Static Assets
    participant M as MCP server (in-Worker)
    participant I as content-index.json

    C->>W: GET /some/path
    alt path starts with /mcp
        W->>M: route to MCP server
        M->>I: query content index (in-memory)
        I-->>M: matching entries
        M-->>W: tool result (text-first content + optional UI resource)
        W-->>C: MCP response
    else Accept: text/markdown, or .md URL sibling
        W->>I: look up canonical Markdown for path
        I-->>W: Markdown source (+ optional in-Worker rendering)
        W-->>C: 200 text/markdown (Vary: Accept, User-Agent)
    else default — full site
        W->>A: fetch matching static asset
        A-->>W: prerendered HTML/CSS/JS
        W-->>C: 200 text/html
    end
```

This single branch point in the `fetch` handler is what lets `biblio` serve three audiences (browsers, MCP clients, terminal/agent clients) from one Worker, one deploy, one origin — see [001 → terminal/CLI-optimized view](001-static-site-platform-with-mcp.md#architecture-patterns) for the content-negotiation research behind the middle branch.

## 4. MCP tool-call flow (e.g., `list_entries` on a project-log site)

```mermaid
sequenceDiagram
    participant Host as MCP host (AI assistant)
    participant MCP as In-Worker MCP server
    participant DO as Durable Object (session)
    participant Idx as content-index.json
    participant UI as MCP App UI bundle

    Host->>MCP: tools/call list_entries { tags: ["biblio"] }
    MCP->>DO: load/establish session state
    MCP->>Idx: filter entries by tag (in-memory query)
    Idx-->>MCP: matching entries
    MCP-->>Host: content: Markdown table / ASCII summary (text-first baseline)
    opt Host supports MCP App UI resources
        MCP-->>Host: + ui:// resource reference (e.g. timeline chart)
        Host->>UI: render sandboxed React bundle (postMessage bridge)
        UI-->>Host: rendered, themed via the site's shared tokens.css
    end
```

The text-first `content` is always returned — terminal-based hosts stop there; graphical hosts additionally get the richer UI resource. See [001 → MCP App UI rendering](001-static-site-platform-with-mcp.md#architecture-patterns) for the confirmed rationale.

The planned `quarterly` type ([§11](roadmap.md#11-quarterly-site-type)) follows the same flow with its own tools (`list_issues`, `get_issue`, `list_articles`, `get_article`). `get_article` returns footnotes and margin notes as structured lists alongside the text, read from the content index rather than parsed per request (pending the §0.5 structured-notes decision).

## 5. Patterns considered (and dropped)

For traceability — these were evaluated against the selected pattern (§1) and dropped; see the [Decision Log](001-static-site-platform-with-mcp.md#decision-log) for the full rationale.

```mermaid
flowchart TB
    subgraph P2["Pattern 2 — single Cloud Run service (dropped)"]
        direction LR
        v2["Visitor"] --> cf2["Cloudflare DNS/proxy<br/>(pass-through only)"] --> cr2["Cloud Run container<br/>serves site + /mcp"]
    end

    subgraph P3["Pattern 3 — Worker as front door, proxy to Cloud Run (dropped)"]
        direction LR
        v3["Visitor"] --> w3["Cloudflare Worker<br/>serves site directly"] -->|"/mcp* proxied via fetch()"| cr3["Cloud Run service<br/>implements MCP"]
    end
```

- **Pattern 2** moves static-asset serving off Cloudflare's free edge onto paid Cloud Run compute, reducing Cloudflare to a thin pass-through for domain/TLS/caching only.
- **Pattern 3** keeps "both/either" backend flexibility but doubles the deployables-per-site (Worker *and* Cloud Run service) and adds proxy/streaming-transparency complexity — more moving parts than `biblio`'s personal/portfolio scale calls for.

## 6. Route vs. subdomain for `/mcp` (within the selected pattern)

```mermaid
flowchart LR
    subgraph Route["Selected: route-based (example.com/mcp)"]
        direction TB
        r1["One Worker"] --> r2["One wrangler.jsonc"] --> r3["One route/domain binding<br/>— no new DNS"]
    end
    subgraph Sub["Rejected: subdomain (mcp.example.com)"]
        direction TB
        s1["New proxied DNS record per zone"] --> s2["Second Worker, or<br/>Host-header branching"] --> s3["More wrangler config<br/>+ DNS to provision & verify, per site"]
    end
```

The route-based split needs no extra DNS work and keeps `biblio-cli` scaffolding to one project/one deploy — strictly simpler to replicate across many sites. CORS, the original concern that motivated considering a subdomain, doesn't actually apply: MCP clients connect directly as external HTTP/SSE clients rather than through the site's own browser JS.

## References

- [001-static-site-platform-with-mcp.md](001-static-site-platform-with-mcp.md) — research and decisions behind every diagram on this page
- [deployment-guide.md](deployment-guide.md) — the end-to-end guide to creating, configuring, and deploying a site on this architecture
- [roadmap.md](roadmap.md) — the implementation task breakdown that builds this architecture out
- [004-quarterly-site-type.md](004-quarterly-site-type.md) — the pre-build step (§2 step 0) and the fourth site type

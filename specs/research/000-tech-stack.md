# Research: Available tech stack for this project

**Status:** Complete
**Date:** 2026-06-07
**Owner:** admin@vixandel.com

## Problem Statement

Establish which of the already-paid-for platforms/services should be used — and how they combine — to build the `biblio` static-site + MCP platform described in [001-static-site-platform-with-mcp.md](001-static-site-platform-with-mcp.md).

- **Goal:** Map each candidate service to a concrete responsibility (hosting, build/deploy, data/auth, AI/LLM, RAG, CI) for `biblio`.
- **In scope:** **Cloudflare Workers** (Pro, 1 domain + many custom domains) and **GitHub** — the two services with assigned, load-bearing responsibilities; **OpenRouter** and **Deepseek** as conditional candidates for build-time embedding generation, contingent on the Vectorize semantic-search escalation in [001-static-site-platform-with-mcp.md](001-static-site-platform-with-mcp.md) actually being built.
- **Out of scope:** Adopting new paid services beyond this list. **GCP** (Cloud Run, Vertex, committed-use Kubernetes discounts), **Vercel**, and **Ragflow** were evaluated against `biblio`'s architecture and found to have no assigned responsibility — see Decision Log for the per-service rationale.
- **Success criteria:** A documented platform → responsibility mapping that `biblio` can be built on without incurring significant new cost.
- **Non-goals:** Cost optimization beyond what's already committed; using every listed service (the project does not need to use all of them); evaluating services outside this list.

## Background & Context

- **Relevant prior work:** [001-static-site-platform-with-mcp.md](001-static-site-platform-with-mcp.md) defines the target system this stack must support.
- **Domain assumptions:** All listed services are already being paid for; incremental usage is expected to add little additional cost; the project should draw from this list rather than introduce new paid dependencies.
- **Stakeholders:** Project owner/admin (admin@vixandel.com), sole developer and user.

## Libraries & Stack

| Candidate | Purpose | License | Maturity | Maintenance | Fit | Notes |
| --- | --- | --- | --- | --- | --- | --- |
| Cloudflare Workers | Edge hosting for the built site *and* the in-Worker `/mcp` endpoint; static assets, MCP session state (Durable Objects), and storage on escalation (D1, Vectorize, KV) | Commercial (Pro, 1 domain) | Mature | Active | **Selected** | Carries essentially all of `biblio`'s compute/hosting/storage responsibilities — see [001-static-site-platform-with-mcp.md](001-static-site-platform-with-mcp.md) for the selected in-Worker architecture |
| GitHub | Source control, CI/CD (Actions) | Commercial | Mature | Active | **Selected** | Drives the build/deploy pipeline: runs `biblio-cli`'s content-index generation, `react-router build` (native prerendering), and `wrangler deploy` on push |
| OpenRouter | Unified API/router across multiple LLM providers, incl. a dedicated embeddings endpoint | Commercial | Mature | Active | **Conditional — preferred** | Candidate for generating build-time embeddings, but only if/when the Vectorize semantic-search escalation in [001](001-static-site-platform-with-mcp.md) is actually built — a thin, build-time-only dependency, not infrastructure. **Researched 2026-06-07:** OpenRouter ships a unified, OpenAI-compatible `/embeddings` endpoint giving access to many embedding models in one place (OpenAI `text-embedding-3-small`/`-large`, Google Gemini Embedding, BGE-M3, Qwen3 Embedding, etc.) — exactly the "thin, build-time-only" fit `biblio` needs |
| Deepseek | LLM provider (chat/reasoning only) | Commercial | Mature | Active | **Ruled out for embeddings** | **Researched 2026-06-07:** confirmed via the official API docs (`api-docs.deepseek.com`) that DeepSeek's API offers **no embeddings endpoint at all** — only chat/reasoning models (`deepseek-v4-flash`/`-pro`). Open GitHub issues on the `deepseek-ai` repos requesting an "embedding-aware API" corroborate this gap. Not a viable candidate for the conditional embedding-generation role; remains available only as a general LLM provider if some other (non-embedding) need arises |

- **Selected:** **Cloudflare Workers** (hosting, MCP, storage) and **GitHub** (source control, CI/CD) are the only load-bearing services `biblio` depends on.
- **Rejected (with reason):**
  - **Supabase** — `biblio` will build and host its own MCP server rather than rely on any third-party/managed MCP or BaaS offering (decided 2026-06-07; see [001-static-site-platform-with-mcp.md](001-static-site-platform-with-mcp.md) Decision Log).
  - **GCP** (Cloud Run, Vertex, committed-use Kubernetes discounts) — Cloud Run was evaluated as an MCP-hosting backend and deprioritized once the in-Worker pattern was selected (see [001](001-static-site-platform-with-mcp.md) Decision Log); nothing in `biblio`'s scope calls for Vertex or Kubernetes-shaped workloads. Stays "available but unused."
  - **Vercel** — already ruled out for MCP hosting by the same-domain/no-CORS decision; nothing else in `biblio`'s scope needs a separate frontend host.
  - **Ragflow** — a full RAG platform; `biblio`'s sketched semantic-search approach (build-time embeddings → Vectorize → in-Worker similarity query) is far lighter-weight and already covers the one plausible use case (the docs site's `search_docs`), so adopting Ragflow would be more infrastructure than the problem calls for.
- **Open questions:** none remaining for this stack-selection research — the provider choice (OpenRouter over DeepSeek) is resolved below; the one narrower follow-on (which embedding *model* within OpenRouter) is tracked in [roadmap.md §0](roadmap.md#0-outstanding-design-questions).

## Architecture Patterns

- **Candidate patterns:** Cloudflare-centric — a single Cloudflare Worker serves the built site and implements the `/mcp` endpoint in-Worker (Workers Static Assets, Durable Objects, escalating to D1/Vectorize as needed), with GitHub driving CI/CD. See [001-static-site-platform-with-mcp.md](001-static-site-platform-with-mcp.md) for the detailed MCP-hosting options analysis and the full reasoning behind selecting this pattern over GCP Cloud Run-based alternatives (which were considered and dropped).
- **Selected pattern(s):** Cloudflare-centric — single in-Worker Cloudflare Worker for site + `/mcp`, GitHub for CI/CD (matches the selected Pattern 1 in [001](001-static-site-platform-with-mcp.md)).
- **Tradeoffs:** Concentrating on Cloudflare + GitHub trades away any use of already-committed GCP resources (Cloud Run, Vertex, Kubernetes discounts) in exchange for the simplest possible topology — one project, one deploy, one origin, no second backend to operate. Given `biblio`'s personal/portfolio scale, operational simplicity wins; the GCP spend remains sunk but unused for this project.
- **Integration points:** Cloudflare Workers ↔ GitHub (CI/CD/deploys — runs the `biblio-cli` content-index generation, `react-router build` with native prerendering, and `wrangler deploy` pipeline); MCP server ↔ OpenRouter/Deepseek (conditional — only if/when build-time embedding generation for Vectorize-backed semantic search is built).
- **Data model implications:** The MCP server's storage lives entirely within Cloudflare (bundled JSON content index by default, escalating to D1 then Vectorize — see [001](001-static-site-platform-with-mcp.md) Decision Log for the full progressive-enhancement approach); no GCP-managed data services are involved.

## Existing Codebase Evaluation

- **Relevant modules:** None — repo is an early-stage scaffold with no application source code yet (per `CLAUDE.md`).
- **Reuse opportunities:** N/A at this stage.
- **Refactor candidates:** N/A
- **Replacements required:** N/A
- **Code health signals (lint, type, test, complexity):** N/A — nothing to measure yet.
- **Conventions to follow:** TBD — to be captured in `AGENTS.md` once stack decisions are made.

## Security & Compliance

- **Threat surface:** Cross-service credentials/secrets (Cloudflare, and conditionally OpenRouter/Deepseek if embedding generation is built); public, unauthenticated `/mcp` endpoint on custom domains (see [001](001-static-site-platform-with-mcp.md) for the public-by-design rationale).
- **Data classification:** Likely low-sensitivity (personal/portfolio use) but should be confirmed per service usage.
- **Compliance regimes:** None anticipated (private, non-commercial use).
- **Mitigations:** Secrets management via Cloudflare's built-in secrets/bindings (no GCP Secret Manager needed, since GCP is not part of the selected stack); Cloudflare's standard protections (rate limiting, WAF, bot management) cover the public `/mcp` endpoint.

## Performance, Reliability & Operations

- **Expected load:** Low — personal/portfolio-scale traffic.
- **Observability:** Cloudflare Workers Observability.
- **Rollout / migration plan:** Incremental — build the scaffold + CLI, validate with one pilot site, then replicate across site types (see [001](001-static-site-platform-with-mcp.md) rollout plan).

## Risks & Mitigations

| Risk | Likelihood | Impact | Mitigation | Owner |
| --- | --- | --- | --- | --- |
| Concentrating on Cloudflare + GitHub creates a single-vendor dependency for hosting/compute | Low | Low | Acceptable given personal/portfolio scale and the operational-simplicity payoff; revisit only if Cloudflare-specific limits become a real constraint | admin@vixandel.com |
| Service choice locks in assumptions that don't fit later site types (project log, docs) | Low | Medium | Validate stack choice against all three intended site types before broad rollout | admin@vixandel.com |

## Open Questions

- [x] ~~If the Vectorize semantic-search escalation in [001](001-static-site-platform-with-mcp.md) is ever built, should it use OpenRouter or Deepseek for build-time embedding generation?~~ — **Resolved by research 2026-06-07: OpenRouter, decisively.** DeepSeek's official API has no embeddings endpoint at all (chat/reasoning models only — confirmed via `api-docs.deepseek.com`); OpenRouter ships a unified, OpenAI-compatible embeddings endpoint spanning many providers' models. The narrower follow-on question — *which* embedding model to use within OpenRouter — is deferred and now tracked in [roadmap.md §0](roadmap.md#0-outstanding-design-questions), since it only matters if/when the Vectorize escalation is actually built — admin@vixandel.com

## Decision Log

- **2026-06-07** — Drop Supabase from the candidate stack; `biblio` will build and host its own MCP server rather than use any third-party/managed MCP or BaaS offering — Owner's stated decision (see [001-static-site-platform-with-mcp.md](001-static-site-platform-with-mcp.md) Decision Log) — admin@vixandel.com
- **2026-06-07** — Finalize the platform → responsibility mapping: **Cloudflare Workers and GitHub are the only load-bearing services** (hosting/MCP/storage and source-control/CI respectively); **OpenRouter and Deepseek remain conditional candidates** for build-time embedding generation, contingent on the Vectorize semantic-search escalation actually being built; **GCP (Cloud Run, Vertex, committed-use Kubernetes), Vercel, and Ragflow are dropped from active consideration** — Owner confirmed this mapping is correct. Rationale: Cloud Run/Vertex were evaluated and deprioritized as an MCP-hosting backend with no AI/ML need identified for `biblio`; Vercel was already ruled out by the same-domain/no-CORS decision; Ragflow is a full RAG platform, heavier than the lightweight build-time-embeddings → Vectorize → in-Worker-similarity-query approach already sketched in [001](001-static-site-platform-with-mcp.md), which already covers the one plausible semantic-search use case — admin@vixandel.com

- **2026-06-07** — Resolve the conditional embedding-provider question in favor of **OpenRouter over DeepSeek**, should the Vectorize semantic-search escalation in [001-static-site-platform-with-mcp.md](001-static-site-platform-with-mcp.md) ever be built — Resolved by direct research rather than a stated preference: DeepSeek's official API (confirmed via `api-docs.deepseek.com`) offers **no embeddings endpoint at all** — chat/reasoning models only — corroborated by open feature-request issues on the `deepseek-ai` GitHub repos asking for one. OpenRouter, by contrast, ships a unified, OpenAI-compatible `/embeddings` endpoint spanning many providers' models (OpenAI, Google Gemini, BGE-M3, Qwen3, etc.), which is exactly the thin, build-time-only dependency `biblio` would need. DeepSeek remains listed only as a general (non-embedding) LLM candidate, now effectively moot since `biblio` has no other identified LLM use — admin@vixandel.com

## References

- [001-static-site-platform-with-mcp.md](001-static-site-platform-with-mcp.md) — defines the system this stack must support

## Appendix

(none yet)

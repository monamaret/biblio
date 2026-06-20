// Tier-ii: in-Worker Markdown→ANSI rendering via `md4x`'s `renderToAnsi()`,
// wired through Wrangler's native `.wasm` import support rather than
// md4x's own auto-detecting `init()` (which only recognizes Node or
// bundler-specific import styles, neither of which Wrangler matches) —
// research.md "In-Worker Markdown→ANSI rendering".
//
// `md4x/wasm`'s own `init()` (the package's public entry point) contains
// an internal dead-code fallback branch with a dynamic `import(".../
// md4x.wasm?url")` for environments where neither an explicit `wasm`
// option nor Node's `fs/promises` is available. That branch is never
// reached here (we always pass an explicit `wasm` module), but esbuild
// statically discovers and tries to bundle every `import()` call site
// regardless of runtime reachability — and Wrangler's bundler has no
// loader for that `?url`-suffixed specifier, breaking the build entirely.
// Importing `common.mjs` directly (the same internals `md4x/wasm`'s
// `init()` itself uses) avoids ever pulling in that broken branch.
import wasmModule from "md4x/build/md4x.wasm";
// A relative filesystem path, not a bare "md4x/..." specifier — md4x's own
// "exports" map doesn't list this subpath, but a literal relative import
// bypasses package-specifier resolution (and its exports-map enforcement)
// entirely, which is exactly how this avoids the broken `init()` above.
import {
  _setInstance,
  _hasInstance,
  _imports,
  renderToAnsi as md4xRenderToAnsi,
  // @ts-expect-error -- deep relative import into md4x's internals; no .d.ts is published for this path.
} from "../../node_modules/md4x/lib/wasm/common.mjs";

let initialized: Promise<void> | undefined;

async function ensureInitialized(): Promise<void> {
  if (_hasInstance()) return;
  if (!initialized) {
    initialized = WebAssembly.instantiate(wasmModule, _imports).then((instance) => {
      _setInstance(instance);
    });
  }
  await initialized;
}

// renderToAnsi renders a Markdown source string to a lightly-formatted
// ANSI string for terminal-capable clients (tier ii). Lazily initializes
// the WASM module exactly once per Worker isolate.
export async function renderToAnsi(markdown: string): Promise<string> {
  await ensureInitialized();
  return md4xRenderToAnsi(markdown);
}

import { afterEach, vi } from "vitest";
import { cleanup } from "@testing-library/react";

// md4x's `.wasm` import (app/negotiation/ansi.ts) only resolves under
// Wrangler's native bundler, not plain Vite/vitest — mock it globally so
// any test that transitively imports ansi.ts (e.g. via accept.ts) doesn't
// need its own boilerplate. ansi.test.ts overrides this with its own
// `vi.mock` + `vi.resetModules()` to assert the wrapper's own logic; tests
// that merely need ansi.ts to be importable (not its actual rendering)
// get a harmless passthrough stub here.
vi.mock("md4x/build/md4x.wasm", () => ({ default: "test-wasm-module" }));

// @testing-library/react doesn't auto-unmount between tests under vitest —
// without this, multiple `it()` blocks that each call `render()` (e.g.
// home.test.tsx's filter-control tests) accumulate DOM across tests,
// producing false "multiple elements found" failures and silently
// over-counting matches in tests that don't fail loudly.
afterEach(() => {
  cleanup();
});

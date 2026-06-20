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

// Carries forward the v0.04 fix for `templates/site/log/vitest.setup.ts`
// (Phase 4 note): without this, multiple `render()` calls across tests in
// one file leak DOM nodes between them, producing false "multiple elements
// found" failures once a test file has more than one `render()` call.
afterEach(() => {
  cleanup();
});

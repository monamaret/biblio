import { vi } from "vitest";

// md4x's `.wasm` import (app/negotiation/ansi.ts) only resolves under
// Wrangler's native bundler, not plain Vite/vitest — mock it globally so
// any test that transitively imports ansi.ts (e.g. via accept.ts) doesn't
// need its own boilerplate. ansi.test.ts overrides this with its own
// `vi.mock` + `vi.resetModules()` to assert the wrapper's own logic; tests
// that merely need ansi.ts to be importable (not its actual rendering)
// get a harmless passthrough stub here.
vi.mock("md4x/build/md4x.wasm", () => ({ default: "test-wasm-module" }));

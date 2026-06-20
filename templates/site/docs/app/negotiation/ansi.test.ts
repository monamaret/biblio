import { describe, expect, it, vi, beforeEach } from "vitest";

// ansi.ts's `.wasm` import only resolves under Wrangler's native bundler
// (confirmed working there per research.md's resolved §0.2 spike) — plain
// Vite/vitest has no WASM loader configured for it, since this file is
// never part of the Vite-bundled React Router app, only the
// Wrangler-bundled Worker entry script. So this test mocks the wasm
// binary import and md4x's low-level instantiation primitives to exercise
// ansi.ts's own lazy-init-once wrapping logic, rather than requiring real
// WASM execution under vitest.
vi.mock("md4x/build/md4x.wasm", () => ({ default: "fake-wasm-module" }));

const renderToAnsiMock = vi.fn((markdown: string) => `ANSI(${markdown})`);
const setInstanceMock = vi.fn();
let hasInstance = false;
vi.mock("../../node_modules/md4x/lib/wasm/common.mjs", () => ({
  _setInstance: setInstanceMock,
  _hasInstance: () => hasInstance,
  _imports: { wasi_snapshot_preview1: {} },
  renderToAnsi: renderToAnsiMock,
}));

const instantiateMock = vi.fn(async () => "fake-instance" as unknown as WebAssembly.Instance);

describe("renderToAnsi wrapper", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.resetModules();
    hasInstance = false;
    vi.stubGlobal("WebAssembly", { instantiate: instantiateMock });
  });

  it("instantiates the WASM module with the imported wasm binding", async () => {
    const { renderToAnsi } = await import("./ansi");
    await renderToAnsi("# Heading");
    expect(instantiateMock).toHaveBeenCalledWith("fake-wasm-module", { wasi_snapshot_preview1: {} });
    expect(setInstanceMock).toHaveBeenCalledWith("fake-instance");
  });

  it("delegates rendering to md4x's renderToAnsi", async () => {
    const { renderToAnsi } = await import("./ansi");
    const result = await renderToAnsi("**bold**");
    expect(renderToAnsiMock).toHaveBeenCalledWith("**bold**");
    expect(result).toBe("ANSI(**bold**)");
  });

  it("only instantiates once across multiple calls", async () => {
    const { renderToAnsi } = await import("./ansi");
    await renderToAnsi("first");
    await renderToAnsi("second");
    expect(instantiateMock).toHaveBeenCalledTimes(1);
  });

  it("skips instantiation entirely if an instance already exists", async () => {
    hasInstance = true;
    const { renderToAnsi } = await import("./ansi");
    await renderToAnsi("already initialized");
    expect(instantiateMock).not.toHaveBeenCalled();
  });
});

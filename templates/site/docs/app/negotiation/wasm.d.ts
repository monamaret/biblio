// Wrangler's native `.wasm` import support produces a WebAssembly.Module
// directly — no bundler-specific loader needed (research.md "In-Worker
// Markdown→ANSI rendering").
declare module "*.wasm" {
  const wasmModule: WebAssembly.Module;
  export default wasmModule;
}

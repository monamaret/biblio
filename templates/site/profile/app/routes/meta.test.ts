import { describe, expect, it } from "vitest";
import { meta as rootMeta } from "../root";
import { meta } from "./detail";
import { pageTitle } from "../content/title";

type Args = Parameters<typeof meta>[0];
const args = (splat: string) => ({ params: { "*": splat } }) as unknown as Args;

// Every page carries its own <title> from content-index.json, never a
// hard-coded string (#215).
describe("page titles", () => {
  it("titles an entry page with the entry's title", () => {
    expect(meta(args("about"))).toEqual([{ title: pageTitle("About Me") }]);
  });

  it("ignores the trailing slash Cloudflare's static-asset serving adds", () => {
    expect(meta(args("about/"))).toEqual([{ title: pageTitle("About Me") }]);
  });

  it("titles an unknown slug as not found", () => {
    expect(meta(args("no/such/page"))).toEqual([{ title: pageTitle("Not found") }]);
  });

  it("falls back to the site name for routes without their own title", () => {
    expect(rootMeta({} as Parameters<typeof rootMeta>[0])).toEqual([{ title: pageTitle() }]);
  });
});

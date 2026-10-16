import { describe, expect, it } from "vitest";
import { meta as rootMeta } from "../root";
import { meta as contentsMeta } from "./contents";
import { meta as archiveMeta } from "./archive";
import { meta as categoryMeta } from "./category";
import { meta as detailMeta } from "./detail";

// Every page carries its own <title> from the index's site name (here the
// template's demo "The Quarterly"), never a hard-coded string (#215).
const args = (params: Record<string, string>) => ({ params }) as unknown as Parameters<typeof detailMeta>[0];

describe("page titles", () => {
  it("titles / with the site name alone", () => {
    expect(contentsMeta(args({}))).toEqual([{ title: "The Quarterly" }]);
    expect(rootMeta(args({}))).toEqual([{ title: "The Quarterly" }]);
  });

  it("titles an issue page with the issue's title", () => {
    expect(contentsMeta(args({ issue: "02-autumn-2026" }))).toEqual([{ title: "Notes from the Margin · The Quarterly" }]);
    expect(contentsMeta(args({ issue: "99-nope" }))).toEqual([{ title: "Not found · The Quarterly" }]);
  });

  it("titles the archive and category pages", () => {
    expect(archiveMeta(args({}))).toEqual([{ title: "Issues · The Quarterly" }]);
    expect(categoryMeta(args({ category: "essays" }))).toEqual([{ title: "Essays · The Quarterly" }]);
    expect(categoryMeta(args({ category: "nope" }))).toEqual([{ title: "Not found · The Quarterly" }]);
  });

  it("titles articles and the colophon, never an issue entry", () => {
    expect(detailMeta(args({ "*": "issues/02-autumn-2026/on-margins/" }))).toEqual([{ title: "On Margins · The Quarterly" }]);
    expect(detailMeta(args({ "*": "colophon" }))).toEqual([{ title: "Colophon · The Quarterly" }]);
    expect(detailMeta(args({ "*": "issues/02-autumn-2026" }))).toEqual([{ title: "Not found · The Quarterly" }]);
  });
});

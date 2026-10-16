import { describe, expect, it } from "vitest";
import { remarkMarginNotes, type MdNode } from "./remark-margin-notes";

const text = (value: string): MdNode => ({ type: "text", value });
const para = (...children: MdNode[]): MdNode => ({ type: "root", children: [{ type: "paragraph", children }] });

function run(tree: MdNode, mode: "inline" | "foot" = "inline"): MdNode {
  remarkMarginNotes({ mode })(tree);
  return tree;
}

function paragraph(tree: MdNode): MdNode[] {
  return tree.children![0].children!;
}

// The plugin turns {>> … <<} in text into marginNote nodes, even when the
// note spans emphasis or links, and leaves code and unclosed markers alone
// (specs/004-quarterly-site-type FR-016, FR-020; research.md R2).
describe("remarkMarginNotes", () => {
  it("wraps a note inside one text node", () => {
    const children = paragraph(run(para(text("A.{>> an aside <<} B."))));
    expect(children.map((c) => c.type)).toEqual(["text", "marginNote", "text"]);
    expect(children[0].value).toBe("A.");
    expect(children[1].children).toEqual([text("an aside")]);
    expect(children[1].data?.hName).toBe("margin-note");
    expect(children[1].data?.hProperties).toEqual({ dataIndex: 1 });
    expect(children[2].value).toBe(" B.");
  });

  it("wraps a note that spans inline formatting", () => {
    const emphasis: MdNode = { type: "emphasis", children: [text("Gerstner")] };
    const children = paragraph(run(para(text("A.{>> cf. "), emphasis, text(", 1964 <<} B."))));
    expect(children.map((c) => c.type)).toEqual(["text", "marginNote", "text"]);
    expect(children[1].children).toEqual([text("cf. "), emphasis, text(", 1964")]);
  });

  it("numbers notes in document order", () => {
    const children = paragraph(run(para(text("{>> one <<} and {>> two <<}"))));
    const notes = children.filter((c) => c.type === "marginNote");
    expect(notes.map((n) => n.data?.hProperties)).toEqual([{ dataIndex: 1 }, { dataIndex: 2 }]);
  });

  it("leaves inline code and unclosed markers as they are", () => {
    const code: MdNode = { type: "inlineCode", value: "{>> no <<}" };
    const children = paragraph(run(para(text("A "), code, text(" and {>> never closed."))));
    expect(children.map((c) => c.type)).toEqual(["text", "inlineCode", "text"]);
    expect(children[2].value).toBe(" and {>> never closed.");
  });

  it("in foot mode, leaves a † reference and lists the notes at the end", () => {
    const tree = run(para(text("A.{>> aside <<} B.")), "foot");
    const children = paragraph(tree);
    expect(children[1].type).toBe("marginNoteRef");
    expect(children[1].data?.hProperties).toEqual({ dataIndex: 1 });
    const foot = tree.children![tree.children!.length - 1];
    expect(foot.type).toBe("marginNoteList");
    expect(foot.children![0].children![0].children).toEqual([text("aside")]);
  });
});

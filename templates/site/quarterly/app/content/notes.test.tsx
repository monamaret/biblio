import { describe, expect, it } from "vitest";
import { render } from "@testing-library/react";
import { ContentBody } from "./render";

const body =
  "Grid.[^1] Margin.{>> cf. *Gerstner* <<} Again.[^2]\n\n[^1]: Brno, 1971.\n[^2]: Second note.\n";

// Notes in HTML (contracts/site-surface.md "Article page"; research.md R7):
// with margin notes on, every footnote has an inline sidenote and a foot
// list entry, and every margin note is a † toggle plus its text — all
// without client JS; with margin notes off, both kinds go to the foot.
describe("ContentBody notes", () => {
  it("renders footnotes as sidenotes plus a foot list when margin notes are on", () => {
    const { container } = render(<ContentBody body={body} marginNotes />);
    const sidenotes = container.querySelectorAll(".sidenote");
    expect(sidenotes).toHaveLength(2);
    expect(sidenotes[0].querySelector(".sidenote-number")?.textContent).toBe("1");
    expect(sidenotes[0].textContent).toContain("Brno, 1971.");
    expect(container.querySelectorAll("section.footnotes li")).toHaveLength(2);
    expect(container.querySelectorAll("sup a[data-footnote-ref]")).toHaveLength(2);
  });

  it("renders a margin note as a † toggle (label + checkbox) and its text", () => {
    const { container } = render(<ContentBody body={body} marginNotes />);
    const label = container.querySelector("label.margin-toggle");
    const checkbox = container.querySelector("input.margin-toggle-checkbox") as HTMLInputElement;
    expect(label?.textContent).toBe("†");
    expect(checkbox.type).toBe("checkbox");
    expect(label?.getAttribute("for")).toBe(checkbox.id);
    const note = container.querySelector("span.marginnote");
    expect(note?.textContent).toBe("cf. Gerstner");
    expect(note?.querySelector("em")?.textContent).toBe("Gerstner");
  });

  it("puts both kinds of note at the foot when margin notes are off", () => {
    const { container } = render(<ContentBody body={body} marginNotes={false} />);
    expect(container.querySelector(".sidenote")).toBeNull();
    expect(container.querySelector(".marginnote")).toBeNull();
    expect(container.querySelectorAll("section.footnotes li")).toHaveLength(2);
    const foot = container.querySelectorAll(".marginnotes-foot li");
    expect(foot).toHaveLength(1);
    expect(foot[0].textContent).toContain("cf. Gerstner");
    expect(container.querySelector("sup.marginnote-ref")?.textContent).toBe("†");
  });

  it("never renders raw HTML inside a note", () => {
    const { container } = render(<ContentBody body={"A.{>> <b>bold</b> <<}"} marginNotes />);
    expect(container.querySelector("b")).toBeNull();
  });
});

import { describe, expect, it } from "vitest";
import { render } from "@testing-library/react";
import { ContentBody } from "./render";
import { listArticles } from "./load";

type Note = { id: string; text: string };

// plain strips the inline Markdown a note's text may carry, to compare it
// with rendered text.
function plain(markdown: string): string {
  return markdown
    .replace(/\[([^\]]*)\]\([^)]*\)/g, "$1")
    .replace(/[*_`]/g, "")
    .trim();
}

// The cross-parser contract (research.md R2): the notes the template
// renders — number, order, text — are exactly the notes biblio-cli
// extracted into the index, for every article in the sample content.
describe("rendered notes agree with the index", () => {
  for (const article of listArticles()) {
    it(article.slug, () => {
      const footnotes = (article.footnotes as Note[]) ?? [];
      const marginNotes = (article.margin_notes as Note[]) ?? [];
      const { container } = render(<ContentBody body={article.body} marginNotes />);

      const sidenotes = Array.from(container.querySelectorAll(".sidenote")).map((el) =>
        (el.textContent ?? "").replace(/^\d+\s*/, "").trim(),
      );
      expect(sidenotes).toEqual(footnotes.map((n) => plain(n.text)));

      const margins = Array.from(container.querySelectorAll("span.marginnote")).map((el) => (el.textContent ?? "").trim());
      expect(margins).toEqual(marginNotes.map((n) => plain(n.text)));
    });
  }
});

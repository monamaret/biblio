import ReactMarkdown, { type Components } from "react-markdown";
import remarkGfm from "remark-gfm";
import { typesetting } from "./load";
import { MarginNote, MarginNoteRef, Sidenote } from "./notes";
import { remarkMarginNotes, remarkSidenotes } from "./remark-margin-notes";

// The custom elements the note plugins emit, mapped to their components.
// react-markdown renders no raw HTML (no rehype-raw), so note text can only
// ever be Markdown (specs/004-quarterly-site-type FR-020).
const components = {
  "margin-note": MarginNote,
  "margin-note-ref": MarginNoteRef,
  "side-note": Sidenote,
} as unknown as Components;

// ContentBody renders an entry's Markdown body via react-markdown +
// remark-gfm — the one rendering path every page uses — with footnotes
// doubled as sidenotes and `{>> … <<}` margin notes, or, when
// typesetting.margin_notes is false, both kinds of note at the foot
// (FR-016–FR-018; research.md R7).
export function ContentBody({ body, marginNotes }: { body: string; marginNotes?: boolean }) {
  const inMargin = marginNotes ?? typesetting().margin_notes;
  const plugins = inMargin
    ? [remarkGfm, remarkSidenotes, [remarkMarginNotes, { mode: "inline" }]]
    : [remarkGfm, [remarkMarginNotes, { mode: "foot" }]];
  return (
    <div className={inMargin ? "article-body with-margin-notes" : "article-body"}>
      <ReactMarkdown remarkPlugins={plugins as never} components={components}>
        {body}
      </ReactMarkdown>
    </div>
  );
}

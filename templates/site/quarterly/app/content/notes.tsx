import type { ReactNode } from "react";

type NoteProps = { children?: ReactNode; "data-index"?: number | string; "data-number"?: number | string };

// MarginNote is the CSS-only margin note (research.md R7, the Tufte
// pattern): a "†" label, a hidden checkbox, and the note text. Wide
// screens float the text into the margin and hide the toggle; narrow
// screens hide the text until the checkbox is ticked. Everything is
// phrasing content, so it is valid inside a paragraph.
export function MarginNote({ children, "data-index": index }: NoteProps) {
  const id = `mn-${index}`;
  return (
    <>
      <label htmlFor={id} className="margin-toggle">
        †
      </label>
      <input type="checkbox" id={id} className="margin-toggle-checkbox" />
      <span className="marginnote">{children}</span>
    </>
  );
}

// MarginNoteRef is the "†" left in the text when margin notes are off; it
// links to the note in the list at the foot of the article.
export function MarginNoteRef({ "data-index": index }: NoteProps) {
  return (
    <sup className="marginnote-ref">
      <a href={`#mn-${index}`}>†</a>
    </sup>
  );
}

// Sidenote is a footnote's text shown beside its reference on wide
// screens; narrow screens hide it and use remark-gfm's list at the foot.
export function Sidenote({ children, "data-number": number }: NoteProps) {
  return (
    <span className="sidenote">
      <span className="sidenote-number">{number}</span> {children}
    </span>
  );
}

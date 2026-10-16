// MdNode is the slice of the mdast shape these plugins touch — kept local
// so the template needs no extra type dependency.
export type MdNode = {
  type: string;
  value?: string;
  children?: MdNode[];
  identifier?: string;
  data?: { hName?: string; hProperties?: Record<string, unknown> };
};

const OPEN = "{>>";
const CLOSE = "<<}";

// remarkMarginNotes turns `{>> … <<}` margin notes into nodes the renderer
// maps to the CSS-only † toggle (specs/004-quarterly-site-type FR-016,
// FR-018; research.md R7). A note may span inline formatting (emphasis,
// links), so markers are matched across a paragraph's sibling nodes; code
// is never a text node, so `{>>` inside code is left alone, as is an
// opener with no closer. Biblio-cli extracts the same notes Go-side for the
// index; a template contract test keeps the two in agreement
// (research.md R2).
//
// mode "inline" (margin_notes: true) keeps each note where it is; mode
// "foot" (margin_notes: false) leaves a † reference and lists the notes at
// the end of the body.
export function remarkMarginNotes(options: { mode: "inline" | "foot" }) {
  return (tree: MdNode) => {
    const collected: MdNode[] = [];
    visitParents(tree, (parent) => {
      parent.children = rewrite(parent.children!, collected, options.mode);
    });
    if (options.mode === "foot" && collected.length > 0) {
      tree.children!.push({
        type: "marginNoteList",
        data: { hName: "section", hProperties: { className: ["marginnotes-foot"] } },
        children: [
          {
            type: "marginNoteItems",
            data: { hName: "ol" },
            children: collected.map((note, i) => ({
              type: "marginNoteItem",
              data: { hName: "li", hProperties: { id: `mn-${i + 1}` } },
              children: note.children,
            })),
          },
        ],
      });
    }
  };
}

// visitParents calls fn for every node whose children are phrasing
// content (paragraphs, headings, table cells, …), depth first.
function visitParents(node: MdNode, fn: (parent: MdNode) => void) {
  if (!node.children) return;
  if (node.children.some((c) => c.type === "text")) {
    fn(node);
  }
  for (const child of node.children) {
    visitParents(child, fn);
  }
}

function rewrite(children: MdNode[], collected: MdNode[], mode: "inline" | "foot"): MdNode[] {
  const out: MdNode[] = [];
  const queue = [...children];
  while (queue.length > 0) {
    const node = queue.shift()!;
    const open = node.type === "text" ? node.value!.indexOf(OPEN) : -1;
    if (open < 0) {
      out.push(node);
      continue;
    }

    // Find the closer: in the rest of this text node, or a later sibling.
    const afterOpen = node.value!.slice(open + OPEN.length);
    let inner: MdNode[] = [];
    let rest = "";
    let consumed = 0;
    let closed = false;
    const sameNode = afterOpen.indexOf(CLOSE);
    if (sameNode >= 0) {
      inner = [{ type: "text", value: afterOpen.slice(0, sameNode) }];
      rest = afterOpen.slice(sameNode + CLOSE.length);
      closed = true;
    } else {
      inner = afterOpen ? [{ type: "text", value: afterOpen }] : [];
      for (let k = 0; k < queue.length; k++) {
        const sib = queue[k];
        const close = sib.type === "text" ? sib.value!.indexOf(CLOSE) : -1;
        if (close >= 0) {
          if (close > 0) inner.push({ type: "text", value: sib.value!.slice(0, close) });
          rest = sib.value!.slice(close + CLOSE.length);
          consumed = k + 1;
          closed = true;
          break;
        }
        inner.push(sib);
      }
    }

    if (!closed) {
      // Unclosed: leave the text exactly as written.
      out.push(node);
      continue;
    }

    queue.splice(0, consumed);
    const before = node.value!.slice(0, open);
    if (before) out.push({ type: "text", value: before });
    const index = collected.length + 1;
    const note: MdNode = {
      type: "marginNote",
      data: { hName: "margin-note", hProperties: { dataIndex: index } },
      children: trim(inner),
    };
    collected.push(note);
    out.push(
      mode === "inline"
        ? note
        : { type: "marginNoteRef", data: { hName: "margin-note-ref", hProperties: { dataIndex: index } } },
    );
    // The remainder may hold another note; scan it next.
    if (rest) queue.unshift({ type: "text", value: rest });
  }
  return out;
}

// trim removes the whitespace just inside the markers.
function trim(nodes: MdNode[]): MdNode[] {
  const result = nodes.map((n) => ({ ...n }));
  const first = result[0];
  if (first?.type === "text") first.value = first.value!.replace(/^\s+/, "");
  const last = result[result.length - 1];
  if (last?.type === "text") last.value = last.value!.replace(/\s+$/, "");
  return result.filter((n) => n.type !== "text" || n.value !== "");
}

// remarkSidenotes copies each footnote definition's text into a sidenote
// node right after its first reference, numbered like remark-gfm's
// footnote list, so wide screens can show it in the margin while narrow
// screens use the list at the foot (research.md R7). Not used when
// margin_notes is false.
export function remarkSidenotes() {
  return (tree: MdNode) => {
    const definitions = new Map<string, MdNode[]>();
    for (const node of tree.children ?? []) {
      if (node.type === "footnoteDefinition" && node.identifier) {
        definitions.set(node.identifier, (node.children ?? []).flatMap((c) => c.children ?? []));
      }
    }
    const numbered = new Map<string, number>();
    const insertAfterRefs = (parent: MdNode) => {
      if (!parent.children || parent.type === "footnoteDefinition") return;
      const out: MdNode[] = [];
      for (const child of parent.children) {
        out.push(child);
        if (child.type === "footnoteReference" && child.identifier && definitions.has(child.identifier) && !numbered.has(child.identifier)) {
          numbered.set(child.identifier, numbered.size + 1);
          out.push({
            type: "sidenote",
            data: { hName: "side-note", hProperties: { dataNumber: numbered.get(child.identifier) } },
            children: structuredClone(definitions.get(child.identifier)!),
          });
        } else {
          insertAfterRefs(child);
        }
      }
      parent.children = out;
    };
    insertAfterRefs(tree);
  };
}

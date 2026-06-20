import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";

// ContentBody renders an entry's Markdown body via react-markdown +
// remark-gfm (research.md "Content loading and rendering in the prerendered
// template") — the one rendering path every detail view uses.
export function ContentBody({ body }: { body: string }) {
  return <ReactMarkdown remarkPlugins={[remarkGfm]}>{body}</ReactMarkdown>;
}

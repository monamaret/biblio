import { type RouteConfig, index, route } from "@react-router/dev/routes";

export default [
  index("routes/contents.tsx", { id: "current-contents" }),
  route("issues", "routes/archive.tsx"),
  route("issues/:issue", "routes/contents.tsx", { id: "issue-contents" }),
  route("category/:category", "routes/category.tsx"),
  // A splat (not ":slug") so article slugs like "issues/02-autumn-2026/on-margins"
  // and the colophon resolve to one detail route.
  route("*", "routes/detail.tsx"),
] satisfies RouteConfig;

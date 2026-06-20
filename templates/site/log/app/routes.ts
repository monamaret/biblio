import { type RouteConfig, index, route } from "@react-router/dev/routes";

export default [
  index("routes/home.tsx"),
  // A splat (not ":slug") so multi-segment slugs like "posts/2026-06-01-hello-world"
  // resolve to one detail route instead of needing one per nesting depth.
  route("*", "routes/detail.tsx"),
] satisfies RouteConfig;

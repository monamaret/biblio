import { type RouteConfig, index, route } from "@react-router/dev/routes";

export default [
  index("routes/home.tsx"),
  // A splat (not ":slug") so multi-segment slugs resolve to one detail
  // route instead of needing one per nesting depth.
  route("*", "routes/detail.tsx"),
] satisfies RouteConfig;

import { describe, expect, it } from "vitest";
import { getProfile, listProjects } from "./profile";

// Every profile-type-specific tool's content response must also be a
// clean, self-sufficient Markdown/plain-text rendering (FR-003).
describe("profile-type-specific MCP tools content contract", () => {
  it("get_profile returns the about page's title and body", () => {
    const result = getProfile();
    const text = (result.content[0] as { text: string }).text;
    expect(text).toContain("About Me");
    expect(text).toContain("Welcome to your new profile site");
  });

  it("list_projects returns every entry under projects/", () => {
    const result = listProjects();
    const text = (result.content[0] as { text: string }).text;
    expect(text).toContain("Sample Project");
    expect(text).toContain("/projects/sample-project");
  });
});

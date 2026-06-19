---
name: speckit-github-close-issues
description: Close all open GitHub issues for the current feature by commenting with blob links to the implementing files on main, then closing each issue.
compatibility: Requires spec-kit project structure with .specify/ directory, gh CLI authenticated, and code already merged to main.
metadata:
  author: biblio-project
user-invocable: true
---

# Close GitHub Issues with Code Links

For every open GitHub issue whose title matches a task in `tasks.md`, post a
comment linking to the file(s) that implement it on `main`, then close the
issue. Run this **after** the feature branch has been merged to main.

## Outline

### 1. Resolve feature directory

Run from the repo root:

```sh
SPECIFY_FEATURE_DIRECTORY=specs/releases/0.01-scaffold-and-theme \
  bash .specify/scripts/bash/check-prerequisites.sh --json --require-tasks --include-tasks
```

Parse `FEATURE_DIR` from the JSON output. If the script fails, fall back to
locating `tasks.md` by searching `specs/` for the directory that matches the
current git branch prefix (same logic as `speckit-implement`).

### 2. Read tasks.md

Read `FEATURE_DIR/tasks.md` in full. Build a task map:

- **Key**: task ID string, e.g. `T001`, `T014`
- **Value**: the full task line, which contains the description and any file
  paths mentioned in backtick code spans (e.g. `` `internal/cli/index.go` ``)

Extract file paths from each task line by finding all `` `path/like/this` ``
spans where the content contains `/` or a file extension (`.go`, `.md`,
`.yaml`, `.yml`, `.ts`, `.tsx`, `.json`, `.toml`, `.sh`, `.css`).

### 3. Get repository info

```sh
gh repo view --json nameWithOwner,url -q '[.nameWithOwner, .url] | @tsv'
```

Capture both `nameWithOwner` (for `gh` calls) and `url` (for blob link
construction). Blob links take the form:

```
{url}/blob/main/{relative-file-path}
```

### 4. Verify files exist on main

For each file path extracted in step 2, check that it exists on the `main`
branch:

```sh
git ls-files --error-unmatch {path}
```

Only include paths that exist. If a task mentions a directory (no file
extension), list its contents with `ls {dir}` and link to the directory
instead:

```
{url}/tree/main/{dir}
```

### 5. List open issues

```sh
gh issue list --state open --limit 200 --json number,title
```

For each issue, parse the task ID from the title: the title begins with
`T` followed by digits and a colon, e.g. `T014: Implement scaffold generator`.
Extract the ID as `T014`. Skip issues whose title does not match this pattern.

### 6. Build and post comments, then close

For each matched issue:

a. Look up the task entry in the task map (step 2) using the parsed task ID.

b. Collect the verified file paths for that task. If there are no verified
   file paths, use a fallback comment noting the task was completed (no
   single-file pointer) and still close the issue.

c. Construct the comment body. Use this format:

   ```
   Implemented in [`{filename}`]({blob-url}).
   
   {one-line description from the task line, after the task ID and flags}
   ```

   If multiple files: list each as a separate bullet:

   ```
   Implemented in:
   - [`{filename1}`]({blob-url-1})
   - [`{filename2}`]({blob-url-2})
   
   {description}
   ```

d. Post the comment:

   ```sh
   gh issue comment {number} --body "{comment}"
   ```

e. Close the issue:

   ```sh
   gh issue close {number}
   ```

f. Print a one-line status: `✓ Closed #N (T014: ...)`.

### 7. Report

After processing all issues, print a summary:

```
Closed N issue(s). M issue(s) skipped (title did not match a task ID).
```

## Graceful Degradation

- **gh not authenticated**: fail immediately with a clear message.
- **No open issues**: print "No open issues matched task IDs — nothing to close."
- **File not on main**: omit the path from the comment; note it was not found.
- **tasks.md not found**: fail with a message directing the user to run `/speckit-tasks` first.
- **Issue already closed**: skip silently.

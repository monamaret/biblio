# biblio-cli

A single binary that scaffolds a complete, runnable, themed site project from one guided command.

## Install

**Download a prebuilt binary** (recommended):

Go to [Releases](../../releases) and download the binary for your platform, then move it onto your `PATH`.

**Or install from source** (requires Go 1.22+):

```sh
go install github.com/monamaret/biblio/cmd/biblio-cli@latest
```

## Usage

### Scaffold a new site

```sh
biblio-cli new my-site --type=profile --domain=example.com --theme=profile-default
```

Site types: `profile`, `log`, `docs`. Omit any flag to be prompted interactively.

### Index authored content

```sh
cd my-site
biblio-cli index
```

Reads every `.md` file under `content/`, validates frontmatter, and writes `content-index.json`. Pass an explicit path to target a different content directory.

### Preview locally

```sh
biblio-cli dev
```

Starts `vite dev` and `wrangler dev` together as one supervised pair. Press `Ctrl+C` to stop both.

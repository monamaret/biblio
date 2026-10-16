// Copyright 2026 Mona Maret
// SPDX-License-Identifier: Apache-2.0

// Package templates embeds the scaffold-template asset tree under site/ so
// that internal/scaffold can read it without relying on the filesystem at
// runtime (Go embed directives cannot reach outside their own directory,
// hence this thin shim sitting next to the assets it embeds).
package templates

import "embed"

// SiteFS holds every file under site/, including dotfiles (the "all:"
// prefix), since some site-type directories currently contain only a
// .gitkeep placeholder pending T013's real seed content.
//
//go:embed all:site
var SiteFS embed.FS

// Copyright 2026 Mona Maret
// SPDX-License-Identifier: Apache-2.0

// Package prebuildregression pins the one change the automatic pre-build
// step is allowed to make to profile, log, and docs projects: it rewrites
// content-index.json. Every other file a scaffold produces must be
// byte-identical whether or not the step ran (specs/004-quarterly-site-type
// FR-036, SC-007).
package prebuildregression_test

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/monamaret/biblio/internal/index"
	"github.com/monamaret/biblio/internal/prebuild"
	"github.com/monamaret/biblio/internal/scaffold"
	"github.com/monamaret/biblio/internal/theme"
)

// snapshot returns every file under dir keyed by slash-separated relative
// path.
func snapshot(t *testing.T, dir string) map[string]string {
	t.Helper()
	files := map[string]string{}
	err := filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		data, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(dir, p)
		files[filepath.ToSlash(rel)] = string(data)
		return nil
	})
	if err != nil {
		t.Fatalf("snapshot %s: %v", dir, err)
	}
	return files
}

func scaffoldSite(t *testing.T, siteType string) string {
	t.Helper()
	dir, err := scaffold.Generate(t.TempDir(), scaffold.Options{
		ProjectName: "regression",
		SiteType:    siteType,
		Domain:      "example.com",
		Theme:       theme.PresetsForType(siteType)[0].ID,
		MCPEnabled:  true,
	})
	if err != nil {
		t.Fatalf("scaffold.Generate(%q): %v", siteType, err)
	}
	return dir
}

func TestPrebuild_ExistingTypesChangeOnlyTheIndex(t *testing.T) {
	for _, siteType := range []string{"profile", "log", "docs"} {
		t.Run(siteType, func(t *testing.T) {
			untouched := snapshot(t, scaffoldSite(t, siteType))

			dir := scaffoldSite(t, siteType)
			if _, err := prebuild.Run(dir, prebuild.Build); err != nil {
				t.Fatalf("prebuild.Run: %v", err)
			}
			built := snapshot(t, dir)

			for path, want := range untouched {
				if path == "content-index.json" {
					continue
				}
				got, ok := built[path]
				if !ok {
					t.Errorf("%s: removed by the pre-build step", path)
				} else if got != want {
					t.Errorf("%s: changed by the pre-build step", path)
				}
			}
			for path := range built {
				if _, ok := untouched[path]; !ok && path != "content-index.json" {
					t.Errorf("%s: created by the pre-build step", path)
				}
			}

			idx, errs := index.Generate(filepath.Join(dir, "content"))
			if len(errs) > 0 {
				t.Fatalf("index.Generate: %v", errs)
			}
			want, _ := json.MarshalIndent(idx, "", "  ")
			if built["content-index.json"] != string(want) {
				t.Error("content-index.json does not equal `biblio-cli index` output")
			}
		})
	}
}

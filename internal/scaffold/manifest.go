// Copyright 2026 Mona Maret
// SPDX-License-Identifier: Apache-2.0

package scaffold

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// manifestPath is where Generate records the scaffold baseline a later
// `biblio-cli update`/`upgrade` diffs against (data-model.md "Scaffold
// Region Classification" classification source). It lives outside
// content/ — Principle II reserves content/ entirely for the author.
const manifestPath = ".biblio/scaffold-manifest.json"

// Manifest is the scaffold baseline recorded at `new` time (and rewritten
// after every successful `update`/`upgrade`): the site type, and a content
// hash for every embedded template file outside content/.
type Manifest struct {
	SiteType string            `json:"site_type"`
	Files    map[string]string `json:"files"`
}

// HashBytes returns the hex-encoded SHA-256 of data, the same hash function
// used for every entry in a Manifest — exported so internal/migrate can
// compute comparable hashes for the new-scaffold and deployed-file sides of
// its 3-way diff.
func HashBytes(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

// BuildManifest hashes every embedded template file for siteType outside
// content/, producing the baseline a freshly scaffolded project starts from.
func BuildManifest(siteType string) (Manifest, error) {
	assets, err := Assets(siteType)
	if err != nil {
		return Manifest{}, fmt.Errorf("scaffold: building manifest for %q: %w", siteType, err)
	}

	files := make(map[string]string)
	err = Walk(siteType, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || strings.HasPrefix(path, "content/") {
			return nil
		}
		data, err := fs.ReadFile(assets, path)
		if err != nil {
			return fmt.Errorf("scaffold: reading template %q: %w", path, err)
		}
		files[path] = HashBytes(data)
		return nil
	})
	if err != nil {
		return Manifest{}, err
	}

	return Manifest{SiteType: siteType, Files: files}, nil
}

// WriteManifest writes m to projectDir's manifest path as indented JSON.
func WriteManifest(projectDir string, m Manifest) error {
	dest := filepath.Join(projectDir, filepath.FromSlash(manifestPath))
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return fmt.Errorf("scaffold: creating manifest directory: %w", err)
	}
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("scaffold: encoding manifest: %w", err)
	}
	if err := os.WriteFile(dest, data, 0o644); err != nil {
		return fmt.Errorf("scaffold: writing manifest: %w", err)
	}
	return nil
}

// ReadManifest reads back the manifest written by WriteManifest for the
// project at projectDir.
func ReadManifest(projectDir string) (Manifest, error) {
	src := filepath.Join(projectDir, filepath.FromSlash(manifestPath))
	data, err := os.ReadFile(src)
	if err != nil {
		return Manifest{}, fmt.Errorf("scaffold: reading manifest: %w", err)
	}
	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return Manifest{}, fmt.Errorf("scaffold: decoding manifest: %w", err)
	}
	return m, nil
}

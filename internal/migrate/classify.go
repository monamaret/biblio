// Copyright 2026 Mona Maret
// SPDX-License-Identifier: Apache-2.0

// Package migrate carries an existing deployed Site Project forward to a
// newer scaffold-template version: classifying each region as
// scaffold-tracked, author-customized, or ambiguous, advancing or
// preserving it accordingly, and surfacing any ambiguous region as an
// Update Conflict rather than silently resolving it (FR-001, FR-002).
package migrate

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/monamaret/biblio/internal/scaffold"
)

// Classification is the per-region judgment ClassifySite makes before any
// write occurs (data-model.md "Scaffold Region Classification").
type Classification string

const (
	// ScaffoldTracked regions are unchanged since the baseline manifest
	// and safe for update to advance to the new scaffold version.
	ScaffoldTracked Classification = "scaffold-tracked"
	// AuthorCustomized regions were edited by the author and the scaffold
	// template hasn't changed them since the baseline — left untouched.
	AuthorCustomized Classification = "author-customized"
	// Ambiguous regions changed on both sides (or otherwise can't be
	// resolved automatically) and become an Update Conflict — never
	// silently resolved in either direction (FR-002).
	Ambiguous Classification = "ambiguous"
)

// Region is one classified path within a Site Project, relative to the
// project root, using forward slashes (matching scaffold.Walk).
type Region struct {
	Path           string
	Classification Classification
	// Current and New are the two sides of the diff that matter once a
	// region is classified: the file as currently deployed, and the file
	// in the scaffold-template version being targeted. A nil slice means
	// the file does not exist on that side (e.g. New is nil when a newer
	// scaffold version removes the path). The manifest's baseline hash is
	// only used internally to derive Classification, since the manifest
	// stores hashes rather than full content.
	Current, New []byte
}

// ClassifySite computes the Scaffold Region Classification for every
// region of the project at projectDir against the embedded scaffold
// template for siteType, using the project's recorded manifest as the
// 3-way-diff baseline (contracts/update-and-escalation.md Step 1). It
// performs no writes — classification always happens before any write
// (data-model.md "Scaffold Region Classification" validation rules).
func ClassifySite(projectDir, siteType string) ([]Region, error) {
	manifest, err := scaffold.ReadManifest(projectDir)
	if err != nil {
		return nil, fmt.Errorf("migrate: reading scaffold manifest: %w", err)
	}

	assets, err := scaffold.Assets(siteType)
	if err != nil {
		return nil, fmt.Errorf("migrate: loading scaffold assets for %q: %w", siteType, err)
	}

	newPaths := make(map[string][]byte)
	err = scaffold.Walk(siteType, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || strings.HasPrefix(path, "content/") || scaffold.IsGenerated(path) {
			return nil
		}
		data, err := fs.ReadFile(assets, path)
		if err != nil {
			return fmt.Errorf("migrate: reading new scaffold template %q: %w", path, err)
		}
		newPaths[path] = data
		return nil
	})
	if err != nil {
		return nil, err
	}

	seen := make(map[string]bool, len(manifest.Files)+len(newPaths))
	var regions []Region

	classifyOne := func(path string) (Region, error) {
		current, currentErr := os.ReadFile(filepath.Join(projectDir, filepath.FromSlash(path)))
		currentExists := currentErr == nil
		if currentErr != nil && !os.IsNotExist(currentErr) {
			return Region{}, fmt.Errorf("migrate: reading deployed file %q: %w", path, currentErr)
		}

		baselineHash, hadBaseline := manifest.Files[path]
		newData, hasNew := newPaths[path]

		var currentHash string
		if currentExists {
			currentHash = scaffold.HashBytes(current)
		}
		var newHash string
		if hasNew {
			newHash = scaffold.HashBytes(newData)
		}

		region := Region{Path: path, Current: current, New: newData}

		switch {
		case !hadBaseline && hasNew && !currentExists:
			// New path introduced by the newer scaffold version, and the
			// author has nothing there yet — safe to add.
			region.Classification = ScaffoldTracked
		case !hadBaseline && hasNew && currentExists:
			// Unexpected collision: a baseline-less path that the new
			// scaffold wants and the author also has. Never guess.
			region.Classification = Ambiguous
		case hadBaseline && !hasNew && !currentExists:
			// Scaffold dropped this path and the author never touched it
			// (it's already gone) — nothing left to do, but report it as
			// scaffold-tracked so callers can drop it from the manifest.
			region.Classification = ScaffoldTracked
		case hadBaseline && !hasNew && currentExists && currentHash == baselineHash:
			// Scaffold dropped this path; the author never touched it —
			// safe to remove.
			region.Classification = ScaffoldTracked
		case hadBaseline && !hasNew && currentExists:
			// Scaffold dropped this path, but the author customized it —
			// never silently delete a customization.
			region.Classification = Ambiguous
		case hadBaseline && hasNew && !currentExists:
			// The author (or something) removed a baseline-tracked file
			// the new scaffold still wants — don't silently resurrect it.
			region.Classification = Ambiguous
		case hadBaseline && hasNew && currentHash == baselineHash:
			// Untouched since scaffolding — safe to advance.
			region.Classification = ScaffoldTracked
		case hadBaseline && hasNew && newHash == baselineHash:
			// The scaffold hasn't changed this path since the baseline;
			// only the author has — preserve their customization.
			region.Classification = AuthorCustomized
		case hadBaseline && hasNew && currentHash == newHash:
			// The author's edit already matches the new scaffold version.
			region.Classification = ScaffoldTracked
		default:
			// Both sides changed, differently — surface for a human
			// decision, never auto-resolve (FR-002).
			region.Classification = Ambiguous
		}

		return region, nil
	}

	for path := range manifest.Files {
		seen[path] = true
		if scaffold.IsGenerated(path) {
			// Build output recorded by an older manifest — never a region.
			continue
		}
		region, err := classifyOne(path)
		if err != nil {
			return nil, err
		}
		regions = append(regions, region)
	}
	for path := range newPaths {
		if seen[path] {
			continue
		}
		region, err := classifyOne(path)
		if err != nil {
			return nil, err
		}
		regions = append(regions, region)
	}

	return regions, nil
}

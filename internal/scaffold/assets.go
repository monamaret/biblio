// Package scaffold generates new Site Project directories from the embedded
// template assets under templates/site/.
package scaffold

import (
	"io/fs"

	"github.com/monamaret/biblio/templates"
)

// SiteTypes are the supported scaffold template trees, matching the
// directories embedded under templates/site/.
var SiteTypes = []string{"profile", "log", "docs", "quarterly"}

// Assets returns the embedded template tree for the given site type, rooted
// so that paths inside it are relative to the type's own directory (e.g.
// "content/site.yaml" rather than "site/profile/content/site.yaml").
func Assets(siteType string) (fs.FS, error) {
	return fs.Sub(templates.SiteFS, "site/"+siteType)
}

// skippedDirs are directories never copied from a template tree. go:embed
// all:site embeds whatever is on disk, so a template whose own tests were
// run locally carries an installed node_modules/ (and React Router's
// generated .react-router/) that must not leak into new projects.
var skippedDirs = map[string]bool{"node_modules": true, ".react-router": true}

// Walk visits every real template file for the given site type, skipping
// .gitkeep VCS placeholders that stand in for not-yet-authored directories
// and the skippedDirs above. Other dotfiles (e.g. quarterly's .gitignore)
// are real template files and are visited.
func Walk(siteType string, fn fs.WalkDirFunc) error {
	assets, err := Assets(siteType)
	if err != nil {
		return err
	}
	return fs.WalkDir(assets, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && skippedDirs[d.Name()] {
			return fs.SkipDir
		}
		if !d.IsDir() && d.Name() == ".gitkeep" {
			return nil
		}
		return fn(path, d, err)
	})
}

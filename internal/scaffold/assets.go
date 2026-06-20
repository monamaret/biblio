// Package scaffold generates new Site Project directories from the embedded
// template assets under templates/site/.
package scaffold

import (
	"io/fs"
	"strings"

	"github.com/monamaret/biblio/templates"
)

// SiteTypes are the supported scaffold template trees, matching the
// directories embedded under templates/site/.
var SiteTypes = []string{"profile", "log", "docs"}

// Assets returns the embedded template tree for the given site type, rooted
// so that paths inside it are relative to the type's own directory (e.g.
// "content/site.yaml" rather than "site/profile/content/site.yaml").
func Assets(siteType string) (fs.FS, error) {
	return fs.Sub(templates.SiteFS, "site/"+siteType)
}

// Walk visits every real template file for the given site type, skipping
// VCS placeholders such as .gitkeep that stand in for not-yet-authored
// directories.
func Walk(siteType string, fn fs.WalkDirFunc) error {
	assets, err := Assets(siteType)
	if err != nil {
		return err
	}
	return fs.WalkDir(assets, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasPrefix(d.Name(), ".") {
			return nil
		}
		return fn(path, d, err)
	})
}

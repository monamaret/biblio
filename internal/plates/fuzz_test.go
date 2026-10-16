// Copyright 2026 Mona Maret
// SPDX-License-Identifier: Apache-2.0

package plates

import (
	"os"
	"testing"
)

// FuzzDecode feeds arbitrary bytes to the decode path: it must return an
// error or an image within the pixel cap — never panic, and never decode
// past the cap (specs/004-quarterly-site-type FR-030; research.md R4).
func FuzzDecode(f *testing.F) {
	for _, seed := range []string{
		"../../testdata/quarterly/issues/01-spring-2026/cover.jpg",
		"../../testdata/quarterly/issues/02-summer-2026/img/alpha.png",
		"../../testdata/quarterly/issues/02-summer-2026/img/alpha.webp",
		"../../testdata/quarterly-bad-images/huge-header.png",
		"../../testdata/quarterly-bad-images/corrupt.jpg",
	} {
		data, err := os.ReadFile(seed)
		if err != nil {
			f.Fatal(err)
		}
		f.Add(data)
	}
	const maxPixels = 1_000_000
	f.Fuzz(func(t *testing.T, data []byte) {
		img, _, err := decodeImage(data, maxPixels)
		if err != nil {
			return
		}
		b := img.Bounds()
		if b.Dx()*b.Dy() > maxPixels {
			t.Fatalf("decoded %dx%d past the %d-pixel cap", b.Dx(), b.Dy(), maxPixels)
		}
	})
}

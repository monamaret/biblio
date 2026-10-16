package plates_test

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/monamaret/biblio/internal/plates"
)

const fixture = "../../testdata/quarterly"

var defaults = plates.Settings{Grayscale: true, MaxWidth: 1600, MaxPixels: 40_000_000}

// copyContent copies the fixture site into a temp content/ dir and returns
// it plus an empty output dir.
func copyContent(t *testing.T) (string, string) {
	t.Helper()
	root := t.TempDir()
	content := filepath.Join(root, "content")
	err := filepath.WalkDir(fixture, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(fixture, p)
		target := filepath.Join(content, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		data, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, 0o644)
	})
	if err != nil {
		t.Fatal(err)
	}
	return content, filepath.Join(root, "public", "plates")
}

func decode(t *testing.T, outDir, url string) (image.Image, string) {
	t.Helper()
	f, err := os.Open(filepath.Join(outDir, strings.TrimPrefix(url, "/plates/")))
	if err != nil {
		t.Fatalf("opening plate %s: %v", url, err)
	}
	defer f.Close()
	img, format, err := image.Decode(f)
	if err != nil {
		t.Fatalf("decoding plate %s: %v", url, err)
	}
	return img, format
}

func isGray(img image.Image) bool {
	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y += 3 {
		for x := b.Min.X; x < b.Max.X; x += 3 {
			r, g, bl, _ := img.At(x, y).RGBA()
			if r != g || g != bl {
				return false
			}
		}
	}
	return true
}

var plateName = regexp.MustCompile(`^/plates/[0-9a-f]{16}\.(jpg|png)$`)

// Every image under content/ becomes a grayscale, width-capped plate with
// a hash name; JPEG stays JPEG, PNG stays PNG and keeps its alpha, WebP
// becomes PNG with alpha or JPEG without (specs/004-quarterly-site-type
// FR-028, FR-029; research.md R4).
func TestConvert_FormatsGrayscaleAndWidth(t *testing.T) {
	content, out := copyContent(t)
	urls, stats, err := plates.Convert(content, out, defaults)
	if err != nil {
		t.Fatalf("Convert() error = %v", err)
	}
	if stats.Converted != 5 || stats.Skipped != 0 {
		t.Errorf("stats = %+v, want 5 converted", stats)
	}
	cases := map[string]string{
		"issues/01-spring-2026/cover.jpg":      "jpeg",
		"issues/02-summer-2026/img/color.jpg":  "jpeg",
		"issues/02-summer-2026/img/alpha.png":  "png",
		"issues/02-summer-2026/img/photo.webp": "jpeg",
		"issues/02-summer-2026/img/alpha.webp": "png",
	}
	for src, wantFormat := range cases {
		url, ok := urls[src]
		if !ok {
			t.Errorf("no plate for %s (map: %v)", src, urls)
			continue
		}
		if !plateName.MatchString(url) {
			t.Errorf("%s → %s, want /plates/<16 hex>.<jpg|png>", src, url)
		}
		img, format := decode(t, out, url)
		if format != wantFormat {
			t.Errorf("%s: format %s, want %s", src, format, wantFormat)
		}
		if !isGray(img) {
			t.Errorf("%s: plate is not grayscale", src)
		}
	}

	wide, _ := decode(t, out, urls["issues/02-summer-2026/img/color.jpg"])
	if b := wide.Bounds(); b.Dx() != 1600 || b.Dy() != 36 {
		t.Errorf("color.jpg plate = %dx%d, want 1600x36 (1800x40 scaled, aspect kept)", b.Dx(), b.Dy())
	}
	small, _ := decode(t, out, urls["issues/01-spring-2026/cover.jpg"])
	if b := small.Bounds(); b.Dx() != 48 {
		t.Errorf("cover.jpg plate width = %d, want 48 (never upscaled)", b.Dx())
	}

	alpha, _ := decode(t, out, urls["issues/02-summer-2026/img/alpha.png"])
	if _, _, _, a := alpha.At(0, 0).RGBA(); a != 0 {
		t.Errorf("alpha.png corner alpha = %d, want 0 (transparency kept)", a)
	}
	if _, _, _, a := alpha.At(16, 16).RGBA(); a == 0 {
		t.Error("alpha.png centre is transparent, want the opaque disc")
	}
}

// grayscale_images: false keeps colour (still width-capped), and the
// setting is part of the hash, so the plate name changes.
func TestConvert_ColorWhenGrayscaleOff(t *testing.T) {
	content, out := copyContent(t)
	gray, _, _ := plates.Convert(content, out, defaults)
	color := defaults
	color.Grayscale = false
	urls, _, err := plates.Convert(content, out, color)
	if err != nil {
		t.Fatal(err)
	}
	src := "issues/02-summer-2026/img/color.jpg"
	if urls[src] == gray[src] {
		t.Error("plate name did not change with the grayscale setting")
	}
	img, _ := decode(t, out, urls[src])
	if isGray(img) {
		t.Error("plate is grayscale with grayscale_images: false")
	}
	if img.Bounds().Dx() != 1600 {
		t.Errorf("width = %d, want 1600 (still capped)", img.Bounds().Dx())
	}
}

// Same bytes and settings give the same name, so a rerun converts
// nothing; changing one source converts one; plates whose source is gone
// are removed (FR-029, FR-031; research.md R4).
func TestConvert_CacheAndCleanup(t *testing.T) {
	content, out := copyContent(t)
	first, _, err := plates.Convert(content, out, defaults)
	if err != nil {
		t.Fatal(err)
	}

	second, stats, err := plates.Convert(content, out, defaults)
	if err != nil {
		t.Fatal(err)
	}
	if stats.Converted != 0 || stats.Skipped != 5 || stats.Removed != 0 {
		t.Errorf("rerun stats = %+v, want 0 converted, 5 skipped", stats)
	}
	for k, v := range first {
		if second[k] != v {
			t.Errorf("%s: name changed between runs (%s → %s)", k, v, second[k])
		}
	}

	cover := filepath.Join(content, "issues", "01-spring-2026", "cover.jpg")
	other, _ := os.ReadFile(filepath.Join(content, "issues", "02-summer-2026", "img", "color.jpg"))
	if err := os.WriteFile(cover, other, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(content, "issues", "02-summer-2026", "img", "photo.webp")); err != nil {
		t.Fatal(err)
	}
	_, stats, err = plates.Convert(content, out, defaults)
	if err != nil {
		t.Fatal(err)
	}
	// The new cover has the same bytes as color.jpg, so it shares that plate.
	if stats.Converted != 0 || stats.Skipped != 4 || stats.Removed != 2 {
		t.Errorf("after edits stats = %+v, want 0 converted, 4 skipped, 2 removed", stats)
	}
	entries, _ := os.ReadDir(out)
	if len(entries) != 3 {
		t.Errorf("plates dir has %d files, want 3", len(entries))
	}
}

// Corrupt files and files over the pixel cap fail with a message naming
// the file, the cap checked from the header alone (FR-030).
func TestConvert_RejectsBadImages(t *testing.T) {
	for name, want := range map[string]string{
		"corrupt.jpg":     "corrupt.jpg",
		"huge-header.png": "40 megapixel",
	} {
		t.Run(name, func(t *testing.T) {
			content, out := copyContent(t)
			data, err := os.ReadFile(filepath.Join("../../testdata/quarterly-bad-images", name))
			if err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(content, "issues", "02-summer-2026", "img", name), data, 0o644); err != nil {
				t.Fatal(err)
			}
			_, _, err = plates.Convert(content, out, defaults)
			if err == nil {
				t.Fatal("Convert() error = nil, want a bad-image error")
			}
			for _, sub := range []string{"issues/02-summer-2026/img/" + name, want} {
				if !strings.Contains(err.Error(), sub) {
					t.Errorf("error %q missing %q", err, sub)
				}
			}
		})
	}
}

// An image reached through a symlink that leaves content/ is rejected.
func TestConvert_RejectsSymlinkOutsideContent(t *testing.T) {
	content, out := copyContent(t)
	outside := filepath.Join(t.TempDir(), "secret.png")
	data, _ := os.ReadFile(filepath.Join(content, "issues", "02-summer-2026", "img", "alpha.png"))
	if err := os.WriteFile(outside, data, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(content, "issues", "02-summer-2026", "img", "link.png")); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	_, _, err := plates.Convert(content, out, defaults)
	if err == nil || !strings.Contains(err.Error(), "link.png") || !strings.Contains(err.Error(), "outside content/") {
		t.Errorf("Convert() error = %v, want link.png rejected as outside content/", err)
	}
}

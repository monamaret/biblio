// Copyright 2026 Mona Maret
// SPDX-License-Identifier: Apache-2.0

// Package plates converts a quarterly site's content images into the
// published "plates": grayscale (unless turned off), width-capped copies
// in public/plates/, named by a hash of the source bytes and the
// conversion settings so the name itself is the cache
// (specs/004-quarterly-site-type FR-028–FR-031; research.md R4).
//
// Everything here is pure Go — no cgo — so goreleaser's cross-builds stay
// simple.
package plates

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/draw"
	_ "golang.org/x/image/webp"
)

// Settings are the conversion knobs; every one of them is part of each
// plate's hash.
type Settings struct {
	// Grayscale converts to black and white (typesetting.grayscale_images).
	Grayscale bool
	// MaxWidth caps the published width in pixels; narrower images are
	// never upscaled.
	MaxWidth int
	// MaxPixels is the decompression-bomb guard, checked from the header
	// before any full decode.
	MaxPixels int
}

// DefaultSettings are the plan's limits: 1600 px wide, 40 megapixels.
var DefaultSettings = Settings{Grayscale: true, MaxWidth: 1600, MaxPixels: 40_000_000}

// encoderVersion is bumped whenever the conversion output changes, so old
// plates are not mistaken for current ones.
const encoderVersion = "plates/1"

// Stats counts what one Convert did, for the pre-build report.
type Stats struct {
	Converted, Skipped, Removed int
}

// imageExts are the file extensions treated as content images.
var imageExts = map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}

// Convert publishes every image under contentDir into outDir and returns
// a map from each source's content-relative path to its plate URL
// (/plates/<hash>.<ext>). Plates that already exist are skipped; plates no
// source maps to any more are deleted. Problems — a corrupt file, one over
// the pixel cap, a symlink leaving contentDir — are all collected and
// returned together, each naming its file.
func Convert(contentDir, outDir string, s Settings) (map[string]string, Stats, error) {
	var stats Stats
	root, err := filepath.EvalSymlinks(contentDir)
	if err != nil {
		return nil, stats, fmt.Errorf("plates: resolving content/: %w", err)
	}

	type job struct {
		rel  string
		data []byte
	}
	var jobs []job
	var errs []error
	walkErr := filepath.WalkDir(contentDir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !imageExts[strings.ToLower(filepath.Ext(p))] {
			return nil
		}
		rel, _ := filepath.Rel(contentDir, p)
		rel = filepath.ToSlash(rel)
		real, err := filepath.EvalSymlinks(p)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", rel, err))
			return nil
		}
		if inside, _ := filepath.Rel(root, real); inside == ".." || strings.HasPrefix(inside, ".."+string(filepath.Separator)) {
			errs = append(errs, fmt.Errorf("%s: resolves outside content/ (symlink to %s); images must live inside content/", rel, real))
			return nil
		}
		data, err := os.ReadFile(real)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", rel, err))
			return nil
		}
		jobs = append(jobs, job{rel, data})
		return nil
	})
	if walkErr != nil {
		return nil, stats, fmt.Errorf("plates: walking content/: %w", walkErr)
	}

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return nil, stats, fmt.Errorf("plates: creating %s: %w", outDir, err)
	}
	urls := map[string]string{}
	keep := map[string]bool{}
	for _, j := range jobs {
		name, err := convertOne(j.data, outDir, s, &stats)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", j.rel, err))
			continue
		}
		urls[j.rel] = "/plates/" + name
		keep[name] = true
	}
	if len(errs) > 0 {
		lines := make([]string, len(errs))
		for i, e := range errs {
			lines[i] = "  • " + e.Error()
		}
		return nil, stats, fmt.Errorf("plates: %d image(s) could not be converted:\n%s", len(errs), strings.Join(lines, "\n"))
	}

	existing, _ := os.ReadDir(outDir)
	for _, e := range existing {
		if !e.IsDir() && !keep[e.Name()] {
			if err := os.Remove(filepath.Join(outDir, e.Name())); err == nil {
				stats.Removed++
			}
		}
	}
	return urls, stats, nil
}

// convertOne writes one plate (unless it already exists) and returns its
// file name.
func convertOne(data []byte, outDir string, s Settings, stats *Stats) (string, error) {
	img, format, err := decodeImage(data, s.MaxPixels)
	if err != nil {
		return "", err
	}
	ext := "jpg"
	if format == "png" || (format == "webp" && !opaque(img)) {
		ext = "png"
	}
	sum := sha256.New()
	fmt.Fprintf(sum, "%s|gray=%t|w=%d|", encoderVersion, s.Grayscale, s.MaxWidth)
	sum.Write(data)
	name := hex.EncodeToString(sum.Sum(nil))[:16] + "." + ext

	target := filepath.Join(outDir, name)
	if _, err := os.Stat(target); err == nil {
		stats.Skipped++
		return name, nil
	}

	out := resize(img, s.MaxWidth)
	if s.Grayscale {
		out = grayscale(out)
	}
	var buf bytes.Buffer
	if ext == "png" {
		err = png.Encode(&buf, out)
	} else {
		err = jpeg.Encode(&buf, out, &jpeg.Options{Quality: 82})
	}
	if err != nil {
		return "", fmt.Errorf("encoding plate: %w", err)
	}
	// Write then rename, so an interrupted build never leaves a truncated
	// plate under its final (cache-hit) name.
	tmp := target + ".tmp"
	if err := os.WriteFile(tmp, buf.Bytes(), 0o644); err != nil {
		return "", err
	}
	if err := os.Rename(tmp, target); err != nil {
		return "", err
	}
	stats.Converted++
	return name, nil
}

// decodeImage checks the header against maxPixels before decoding the
// pixels, so an oversized or hostile file is rejected without allocating
// its full bitmap.
func decodeImage(data []byte, maxPixels int) (image.Image, string, error) {
	cfg, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return nil, "", fmt.Errorf("not a readable JPEG, PNG or WebP image: %w", err)
	}
	if cfg.Width <= 0 || cfg.Height <= 0 {
		return nil, "", errors.New("image has no pixels")
	}
	if int64(cfg.Width)*int64(cfg.Height) > int64(maxPixels) {
		return nil, "", fmt.Errorf("image is %dx%d, over the %d megapixel limit", cfg.Width, cfg.Height, maxPixels/1_000_000)
	}
	img, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, "", fmt.Errorf("corrupt %s image: %w", format, err)
	}
	if b := img.Bounds(); b.Dx() != cfg.Width || b.Dy() != cfg.Height {
		return nil, "", errors.New("image size does not match its header")
	}
	return img, format, nil
}

func opaque(img image.Image) bool {
	if o, ok := img.(interface{ Opaque() bool }); ok {
		return o.Opaque()
	}
	return false
}

// resize scales img down to maxWidth with Catmull-Rom, keeping its aspect
// ratio; narrower images are returned unchanged.
func resize(img image.Image, maxWidth int) image.Image {
	b := img.Bounds()
	if b.Dx() <= maxWidth {
		return img
	}
	h := (b.Dy()*maxWidth + b.Dx()/2) / b.Dx() // rounded to nearest
	if h < 1 {
		h = 1
	}
	dst := image.NewNRGBA(image.Rect(0, 0, maxWidth, h))
	draw.CatmullRom.Scale(dst, dst.Bounds(), img, b, draw.Src, nil)
	return dst
}

// grayscale converts by Rec. 709 luminance: a single-channel image.Gray
// for opaque images, an NRGBA with equal channels when there's alpha to
// keep.
func grayscale(img image.Image) image.Image {
	b := img.Bounds()
	if opaque(img) {
		g := image.NewGray(b)
		for y := b.Min.Y; y < b.Max.Y; y++ {
			for x := b.Min.X; x < b.Max.X; x++ {
				r, gg, bl, _ := img.At(x, y).RGBA()
				g.SetGray(x, y, color.Gray{Y: luma(r, gg, bl)})
			}
		}
		return g
	}
	out := image.NewNRGBA(b)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			c := color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA)
			l := luma(uint32(c.R)*0x101, uint32(c.G)*0x101, uint32(c.B)*0x101)
			out.SetNRGBA(x, y, color.NRGBA{R: l, G: l, B: l, A: c.A})
		}
	}
	return out
}

// luma is Rec. 709 luminance of 16-bit channels, as an 8-bit value.
func luma(r, g, b uint32) uint8 {
	return uint8((2126*r + 7152*g + 722*b) / 10000 >> 8)
}

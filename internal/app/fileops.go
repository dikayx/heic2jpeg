package app

import (
	"image"
	"image/jpeg"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/gen2brain/heic"
)

func isHeic(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	return ext == ".heic" || ext == ".heif"
}

func replaceExt(p, ext string) string {
	return strings.TrimSuffix(p, filepath.Ext(p)) + ext
}

func decodeHEIC(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	return img, err
}

func writeJPEG(img image.Image, out string, quality int, dryRun bool) error {
	if dryRun {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
		return err
	}
	tmp := out + ".tmp-" + time.Now().String()
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := jpeg.Encode(f, img, &jpeg.Options{Quality: quality}); err != nil {
		return err
	}
	f.Sync()
	return os.Rename(tmp, out)
}

func copyFile(src, dst string, dryRun bool) error {
	if dryRun {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

func deleteFile(path string, dryRun bool) error {
	if dryRun {
		return nil
	}
	return os.Remove(path)
}

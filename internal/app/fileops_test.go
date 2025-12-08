package app

import (
	"image"
	"image/color"
	"os"
	"path/filepath"
	"testing"
)

func TestIsHeic(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want bool
	}{
		{"heic", "file.HEIC", true},
		{"heif", "something.heif", true},
		{"jpeg", "a.jpg", false},
		{"noext", "file", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isHeic(tt.in); got != tt.want {
				t.Fatalf("isHeic(%s) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

func TestReplaceExt(t *testing.T) {
	r := replaceExt("path/image.heic", ".jpg")
	if r != "path/image.jpg" {
		t.Fatalf("replaceExt returned: %s", r)
	}
}

func TestCopyFile(t *testing.T) {
	tmp := t.TempDir()
	src := filepath.Join(tmp, "src.txt")
	dst := filepath.Join(tmp, "dst.txt")

	os.WriteFile(src, []byte("hello"), 0644)

	if err := copyFile(src, dst, false); err != nil {
		t.Fatalf("copyFile returned err: %v", err)
	}

	b, _ := os.ReadFile(dst)
	if string(b) != "hello" {
		t.Fatalf("copyFile content mismatch: %s", string(b))
	}
}

func TestDeleteFile(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "x.txt")
	os.WriteFile(f, []byte("x"), 0644)

	if err := deleteFile(f, false); err != nil {
		t.Fatalf("deleteFile error: %v", err)
	}

	if _, err := os.Stat(f); !os.IsNotExist(err) {
		t.Fatalf("deleteFile did not delete file")
	}
}

func TestWriteJPEG(t *testing.T) {
	tmp := t.TempDir()
	out := filepath.Join(tmp, "img.jpg")

	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	img.Set(0, 0, color.RGBA{255, 0, 0, 255})

	if err := writeJPEG(img, out, 80, false); err != nil {
		t.Fatalf("writeJPEG error: %v", err)
	}

	if _, err := os.Stat(out); err != nil {
		t.Fatalf("output JPEG missing at %s", out)
	}
}

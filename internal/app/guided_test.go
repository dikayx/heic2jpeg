package app

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWalkFiles(t *testing.T) {
	tmp := t.TempDir()
	os.Mkdir(filepath.Join(tmp, "sub"), 0755)

	os.WriteFile(filepath.Join(tmp, "a.txt"), []byte{}, 0644)
	os.WriteFile(filepath.Join(tmp, "sub", "b.txt"), []byte{}, 0644)

	jobs, err := walkFiles(tmp)
	if err != nil {
		t.Fatalf("walkFiles error: %v", err)
	}

	if len(jobs) != 2 {
		t.Fatalf("walkFiles expected 2 files, got %d", len(jobs))
	}
}

func TestCountHeic(t *testing.T) {
	jobs := []job{
		{"a.heic", "a.heic"},
		{"b.jpg", "b.jpg"},
		{"c.HEIF", "c.HEIF"},
	}

	c := countHeic(jobs)
	if c != 2 {
		t.Fatalf("countHeic expected 2, got %d", c)
	}
}

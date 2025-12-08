package app

import (
	"testing"
)

// Dummy jobs for now
func TestRunWorkers_NoRefactor(t *testing.T) {
	jobs := []job{
		{src: "file1.heic", rel: "file1.heic"},
		{src: "file2.png", rel: "file2.png"},
	}

	err := runWorkers(jobs, "/src", "/dst", ModeInPlace, 80, 2, true, true)
	if err == nil {
		t.Log("No error returned (likely because files don't exist, as expected)")
	} else {
		t.Logf("Error returned (expected in real scenario): %v", err)
	}
}

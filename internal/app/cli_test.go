package app

import "testing"

func TestCLIExclusiveModes(t *testing.T) {
	// invalid: multiple modes
	args := []string{"--inplace", "--convert"}

	// capture stdout by redirecting if desired — here just ensure no panic
	runCLI(args)
}

func TestCLIValidFlags(t *testing.T) {
	args := []string{
		"--inplace",
		"--source", "/tmp",
	}

	runCLI(args) // ensures no crash
}

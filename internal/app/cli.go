package app

import (
	"flag"
	"fmt"
	"runtime"
)

func runCLI(args []string) {
	fs := flag.NewFlagSet("cli", flag.ExitOnError)

	inplace := fs.Bool("inplace", false, "In-place")
	convert := fs.Bool("convert", false, "Convert")
	copyAll := fs.Bool("copy", false, "Copy")

	src := fs.String("source", "", "Source folder")
	dst := fs.String("destination", "", "Destination")

	quality := fs.Int("quality", 90, "Quality")
	dryRun := fs.Bool("dry-run", false, "Dry run")
	deleteOrig := fs.Bool("delete-originals", false, "Delete originals")

	workers := fs.Int("workers", runtime.GOMAXPROCS(0), "Workers")

	fs.Parse(args)

	modeCount := 0
	if *inplace {
		modeCount++
	}
	if *convert {
		modeCount++
	}
	if *copyAll {
		modeCount++
	}
	if modeCount != 1 {
		fmt.Println("Choose one mode flag")
		return
	}

	var mode int
	switch {
	case *inplace:
		mode = ModeInPlace
	case *convert:
		mode = ModeDestFolder
	case *copyAll:
		mode = ModeCopyConvert
	}

	jobs, _ := walkFiles(*src)

	err := runWorkers(jobs, *src, *dst, mode, *quality, *workers, *dryRun, *deleteOrig)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Done.")
}

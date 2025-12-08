package app

import (
	"bytes"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
)

const (
	ModeInPlace     = 1
	ModeDestFolder  = 2
	ModeCopyConvert = 3
)

type job struct {
	src string
	rel string
}

// runWorkers processes jobs concurrently and shows a progress bar
func runWorkers(jobs []job, srcRoot, dstRoot string, mode, quality, workers int, dryRun, deleteOrig bool) error {
	if workers <= 0 {
		workers = 1
	}

	var wg sync.WaitGroup
	jobCh := make(chan job)
	errCh := make(chan error, len(jobs))

	total := len(jobs)
	done := 0
	var mu sync.Mutex // protects done counter

	// Start workers
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobCh {
				if err := processJob(j, srcRoot, dstRoot, mode, quality, dryRun, deleteOrig); err != nil {
					errCh <- err
				}

				// Update progress
				mu.Lock()
				done++
				printProgress(done, total)
				mu.Unlock()
			}
		}()
	}

	// Feed jobs
	for _, j := range jobs {
		jobCh <- j
	}
	close(jobCh)

	wg.Wait()
	close(errCh)

	// Collect errors
	var buf bytes.Buffer
	for e := range errCh {
		buf.WriteString(e.Error() + "\n")
	}
	if buf.Len() > 0 {
		return errors.New(buf.String())
	}

	return nil
}

// processJob handles a single job according to mode
func processJob(j job, srcRoot, dstRoot string, mode, quality int, dryRun, deleteOrig bool) error {
	src := j.src
	rel := j.rel
	isH := isHeic(src)

	var err error
	switch mode {
	case ModeInPlace:
		if isH {
			img, e := decodeHEIC(src)
			if e != nil {
				err = e
				break
			}
			err = writeJPEG(img, replaceExt(src, ".jpg"), quality, dryRun)
			if err == nil && deleteOrig {
				_ = deleteFile(src, dryRun)
			}
		}

	case ModeDestFolder:
		dst := filepath.Join(dstRoot, rel)
		if isH {
			img, e := decodeHEIC(src)
			if e != nil {
				err = e
				break
			}
			err = writeJPEG(img, replaceExt(dst, ".jpg"), quality, dryRun)
		}

	case ModeCopyConvert:
		dst := filepath.Join(dstRoot, rel)
		if isH {
			img, e := decodeHEIC(src)
			if e != nil {
				err = e
				break
			}
			err = writeJPEG(img, replaceExt(dst, ".jpg"), quality, dryRun)
		} else {
			err = copyFile(src, dst, dryRun)
		}
	}

	return err
}

// printProgress shows a terminal progress bar
func printProgress(done, total int) {
	percent := float64(done) / float64(total)
	barSize := 30
	filled := int(percent * float64(barSize))
	bar := "[" + strings.Repeat("█", filled) + strings.Repeat(" ", barSize-filled) + "]"
	fmt.Printf("\r%s %3.0f%% (%d/%d)", bar, percent*100, done, total)
	if done == total {
		fmt.Println()
	}
}

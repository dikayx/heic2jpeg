package app

import (
	"bytes"
	"errors"
	"path/filepath"
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

func runWorkers(jobs []job, srcRoot, dstRoot string, mode, quality, workers int, dryRun, deleteOrig bool) error {
	var wg sync.WaitGroup
	jobCh := make(chan job)
	errCh := make(chan error, len(jobs))

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobCh {
				src := j.src
				rel := j.rel
				isH := isHeic(src)

				var err error
				switch mode {
				case ModeInPlace:
					if isH {
						img, e := decodeHEIC(src)
						if e == nil {
							err = writeJPEG(img, replaceExt(src, ".jpg"), quality, dryRun)
							if err == nil && deleteOrig {
								_ = deleteFile(src, dryRun)
							}
						} else {
							err = e
						}
					}
				case ModeDestFolder, ModeCopyConvert:
					dst := filepath.Join(dstRoot, rel)
					if isH {
						img, e := decodeHEIC(src)
						if e == nil {
							err = writeJPEG(img, replaceExt(dst, ".jpg"), quality, dryRun)
						} else {
							err = e
						}
					} else if mode == ModeCopyConvert {
						err = copyFile(src, dst, dryRun)
					}
				}

				if err != nil {
					errCh <- err
				}
			}
		}()
	}

	for _, j := range jobs {
		jobCh <- j
	}
	close(jobCh)
	wg.Wait()
	close(errCh)

	var buf bytes.Buffer
	for e := range errCh {
		buf.WriteString(e.Error() + "\n")
	}
	if buf.Len() > 0 {
		return errors.New(buf.String())
	}
	return nil
}

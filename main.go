package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	_ "github.com/gen2brain/heic" // HEIC decoder
	"github.com/ncruces/zenity"   // improved dialogs with folder creation support
)

const (
	ModeInPlace     = 1
	ModeDestFolder  = 2
	ModeCopyConvert = 3
)

var (
	modeFlag = flag.Int("mode", 0, "Mode (1=in-place, 2=convert->dest, 3=copy->convert->dest)")
	srcFlag  = flag.String("src", "", "Source folder (optional; will prompt)")
	dstFlag  = flag.String("dst", "", "Destination for mode 2/3 (optional; will prompt)")
	workers  = flag.Int("workers", runtime.GOMAXPROCS(0), "Worker count")
	quiet    = flag.Bool("quiet", false, "Reduce output")
)

func printHelp() {
	fmt.Println(`
HEIC2JPEG — Cross-Platform HEIC Converter
-----------------------------------------

Modes:
  -mode=1   In-place conversion
             • Converts all HEIC files in the selected folder to JPEG
             • JPEGs are placed next to originals
             • Folder structure stays exactly the same

  -mode=2   Convert → Destination folder
             • Converts only HEIC files
             • Keeps original files untouched
             • Preserves full folder structure under destination folder

  -mode=3   Copy + Convert along the way
             • Copies ALL files from source → destination
             • Converts HEIC → JPEG during copy
             • Original HEIC files remain in source
             • Destination folder contains JPEGs and all non-HEIC files

Options:
  -src       Source folder (optional; uses visual folder picker if omitted)
  -dst       Destination folder (required for modes 2 and 3 unless picker used)
  -workers   Number of parallel workers (default: number of CPU cores)
  -quiet     Minimal output (no progress bar or messages)
  -help      Show this help text

Examples:
  heicconv -mode=1
  heicconv -mode=2 -src="/Photos/HEIC" -dst="/Converted/JPEG"
  heicconv -mode=3 -src="~/DCIM" -dst="~/Archive"

Notes:
  • Folder structure is always preserved.
  • No external codecs needed (no ffmpeg required).
  • Native folder chooser allows creating new folders.
`)
}

func welcome() {
	fmt.Println("=========================================================")
	fmt.Println("         HEIC → JPEG CONVERTER (CROSS PLATFORM)")
	fmt.Println("=========================================================")
	fmt.Println("This tool converts HEIC images to JPEG with no external")
	fmt.Println("dependencies. Folder structures are preserved.")
	fmt.Println("---------------------------------------------------------")
	fmt.Println()
}

func pause(msg string) {
	fmt.Printf("%s (press ENTER to continue)", msg)
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	fmt.Println()
}

func isHeic(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	return ext == ".heic" || ext == ".heif"
}

func safeWriteJPEG(img image.Image, outPath string, quality int) error {
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return err
	}
	tmp := outPath + fmt.Sprintf(".tmp-%d", time.Now().UnixNano())

	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := jpeg.Encode(f, img, &jpeg.Options{Quality: quality}); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	f.Sync()

	return os.Rename(tmp, outPath)
}

func decodeHEIC(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f) // discard format string
	if err != nil {
		return nil, err
	}
	return img, nil
}

func copyFile(src, dst string) error {
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

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}

type job struct {
	srcPath string
	relPath string
}

func walkSource(root string) ([]job, error) {
	var jobs []job
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(root, path)
		jobs = append(jobs, job{path, rel})
		return nil
	})
	return jobs, err
}

func progressBar(current, total int) {
	percent := float64(current) / float64(total)
	barSize := 30
	filled := int(percent * float64(barSize))

	bar := "[" + strings.Repeat("█", filled) + strings.Repeat(" ", barSize-filled) + "]"
	fmt.Printf("\r%s %3.0f%% (%d/%d)", bar, percent*100, current, total)
	if current == total {
		fmt.Println()
	}
}

func runWorkers(jobs []job, srcRoot, dstRoot string, mode int, quality int) error {
	var wg sync.WaitGroup
	jobCh := make(chan job)
	errCh := make(chan error, len(jobs))

	total := len(jobs)
	done := 0
	mu := sync.Mutex{}

	// Workers
	for i := 0; i < *workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobCh {
				src := j.srcPath
				rel := j.relPath
				isH := isHeic(src)

				var err error

				switch mode {

				case ModeInPlace:
					if isH {
						img, err := decodeHEIC(src)
						if err == nil {
							err = safeWriteJPEG(img, fileReplaceExt(src, ".jpg"), quality)
						}
					}

				case ModeDestFolder:
					dst := filepath.Join(dstRoot, rel)
					if isH {
						dst = fileReplaceExt(dst, ".jpg")
						img, err2 := decodeHEIC(src)
						if err2 == nil {
							err = safeWriteJPEG(img, dst, quality)
						} else {
							err = err2
						}
					} else {
						err = copyFile(src, dst)
					}

				case ModeCopyConvert:
					dst := filepath.Join(dstRoot, rel)
					if isH {
						dst = fileReplaceExt(dst, ".jpg")
						img, err2 := decodeHEIC(src)
						if err2 == nil {
							err = safeWriteJPEG(img, dst, quality)
						} else {
							err = err2
						}
					} else {
						err = copyFile(src, dst)
					}
				}

				if err != nil {
					errCh <- err
				}

				mu.Lock()
				done++
				progressBar(done, total)
				mu.Unlock()
			}
		}()
	}

	go func() {
		for _, j := range jobs {
			jobCh <- j
		}
		close(jobCh)
	}()

	wg.Wait()
	close(errCh)

	// collect errors
	var buf bytes.Buffer
	for e := range errCh {
		buf.WriteString(e.Error() + "\n")
	}
	if buf.Len() > 0 {
		return errors.New(buf.String())
	}
	return nil
}

func pickFolder(prompt string) (string, error) {
	return zenity.SelectFile(
		zenity.Directory(),
		zenity.Title(prompt),
		zenity.ConfirmOverwrite(),
	)
}

func fileReplaceExt(p, newExt string) string {
	return strings.TrimSuffix(p, filepath.Ext(p)) + newExt
}

func main() {
	flag.Parse()

	// Show help automatically if requested or if mode is unset
	if len(os.Args) == 1 || *modeFlag == 0 {
		printHelp()
		return
	}

	welcome()

	if *modeFlag == 0 {
		fmt.Println("No mode selected. Use -mode=1,2,3")
		return
	}

	pause("Next step: select SOURCE folder")

	var err error
	src := *srcFlag
	if src == "" {
		src, err = pickFolder("Select source folder")
		if err != nil {
			fmt.Println("Cancelled.")
			return
		}
	}

	var dst string
	if *modeFlag == ModeDestFolder || *modeFlag == ModeCopyConvert {
		pause("Next step: select DESTINATION folder")
		dst = *dstFlag
		if dst == "" {
			dst, err = pickFolder("Select destination folder")
			if err != nil {
				fmt.Println("Cancelled.")
				return
			}
		}
	}

	fmt.Println("Gathering files…")
	jobs, err := walkSource(src)
	if err != nil {
		fmt.Println("Error reading source:", err)
		return
	}
	fmt.Printf("Found %d files.\n\n", len(jobs))

	pause("Start processing now")

	err = runWorkers(jobs, src, dst, *modeFlag, 90)
	if err != nil {
		fmt.Println("Errors occurred:\n", err)
		return
	}

	fmt.Println("\nAll done! 🎉")
}

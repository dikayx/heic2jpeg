package app

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ncruces/zenity"
)

func runGuided() {
	welcome()
	reader := bufio.NewReader(os.Stdin)

	// --- Step 1: Choose conversion type ---
	fmt.Println("Choose conversion type:")
	fmt.Println("1) In-place")
	fmt.Println("2) Convert to destination")
	fmt.Println("3) Copy + Convert")
	fmt.Print("Enter choice: ")

	var mode int
	fmt.Fscan(reader, &mode)
	reader.ReadString('\n')

	// --- Step 2: Advanced options ---
	quality := 90
	dryRun := false
	deleteOrig := false

	fmt.Print("Configure advanced options? (y/n): ")
	yn, _ := reader.ReadString('\n')
	yn = strings.TrimSpace(strings.ToLower(yn))

	if yn == "y" {
		fmt.Print("JPEG quality (1-100): ")
		qStr, _ := reader.ReadString('\n')
		if q, err := strconv.Atoi(strings.TrimSpace(qStr)); err == nil {
			quality = q
		}

		fmt.Print("Dry run? (y/n): ")
		dStr, _ := reader.ReadString('\n')
		dryRun = strings.TrimSpace(strings.ToLower(dStr)) == "y"

		fmt.Print("Delete originals? (y/n): ")
		delStr, _ := reader.ReadString('\n')
		deleteOrig = strings.TrimSpace(strings.ToLower(delStr)) == "y"
	}

	// --- Step 3: Pick source folder ---
	pause("Next, you'll pick the SOURCE folder.")
	src, err := zenity.SelectFile(zenity.Directory(), zenity.Title("Select source folder"))
	if err != nil {
		fmt.Println("Cancelled.")
		return
	}

	// --- Step 4: Scan files ---
	jobs, _ := walkFiles(src)
	heicCount := countHeic(jobs)

	fmt.Printf("Gathering files…\nFound %d files (%d HEIC).\n", len(jobs), heicCount)
	if mode == ModeInPlace && heicCount == 0 {
		fmt.Println("No HEIC files found. Exiting.")
		return
	}

	// --- Step 5: Pick destination if needed ---
	var dst string
	if mode == ModeDestFolder || mode == ModeCopyConvert {
		pause("Next, you'll choose the DESTINATION folder.")
		dst, _ = zenity.SelectFile(zenity.Directory(), zenity.Title("Select destination folder"))
	}

	// --- Step 6: Process files ---
	err = runWorkers(jobs, src, dst, mode, quality, 4, dryRun, deleteOrig)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Done.")
}

// pause prints a message and waits for ENTER
func pause(msg string) {
	fmt.Printf("%s (press ENTER to continue)\n", msg)
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func welcome() {
	fmt.Println("=======================================")
	fmt.Println(" HEIC → JPEG Converter")
	fmt.Println("=======================================")
}

func walkFiles(root string) ([]job, error) {
	var jobs []job
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			rel, _ := filepath.Rel(root, path)
			jobs = append(jobs, job{path, rel})
		}
		return nil
	})
	return jobs, nil
}

func countHeic(jobs []job) int {
	count := 0
	for _, j := range jobs {
		if isHeic(j.src) {
			count++
		}
	}
	return count
}

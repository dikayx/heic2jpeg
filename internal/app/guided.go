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

	fmt.Println("Choose conversion type:")
	fmt.Println("1) In-place")
	fmt.Println("2) Convert to destination")
	fmt.Println("3) Copy + Convert")

	fmt.Print("Enter choice: ")
	var mode int
	fmt.Fscan(reader, &mode)
	reader.ReadString('\n')

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

	src, err := zenity.SelectFile(zenity.Directory(), zenity.Title("Select source folder"))
	if err != nil {
		fmt.Println("Cancelled.")
		return
	}

	jobs, _ := walkFiles(src)
	heicCount := countHeic(jobs)

	if mode == ModeInPlace && heicCount == 0 {
		fmt.Println("No HEIC files found.")
		return
	}

	var dst string
	if mode == ModeDestFolder || mode == ModeCopyConvert {
		dst, _ = zenity.SelectFile(zenity.Directory(), zenity.Title("Select destination folder"))
	}

	err = runWorkers(jobs, src, dst, mode, quality, 4, dryRun, deleteOrig)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Done.")
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

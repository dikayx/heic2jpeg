package app

import "fmt"

func printHelp() {
	fmt.Println(`
HEIC2JPEG — Cross-Platform HEIC Converter
=========================================

Usage:
  heic2jpeg                 # guided mode
  heic2jpeg cli [flags]
  heic2jpeg help

Modes (CLI):
  --inplace
  --convert
  --copy

Options:
  --quality N              JPEG quality (1–100)
  --dry-run                Simulate without changing files
  --delete-originals       Delete HEIC files after successful conversion

Required Flags:
  --source PATH
  --destination PATH       (for --convert or --copy)

Other:
  --workers N

Examples:
  heic2jpeg cli --inplace --source "/photos" --quality 85
  heic2jpeg cli --convert --source "/in" --destination "/out" --dry-run
  heic2jpeg cli --copy --source "/in" --destination "/out" --delete-originals`)
}

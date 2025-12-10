# **HEIC2JPEG — Cross-Platform HEIC Converter**

A simple, fast, and dependency-free **HEIC → JPEG converter** for macOS, Windows, and Linux.  

✔ Preserves folder structure  
✔ Converts HEIC/HEIF to JPEG  
✔ Optional copying of all other files  
✔ Guided UI mode (with dialogs!)  
✔ Command-line mode for automation  
✔ Supports dry-runs, deletion of originals, and adjustable quality  
✔ Parallel processing (fast!)

---

## ✨ Features

- **Guided mode** (default)  
  Friendly interactive workflow with step-by-step instructions and system-native folder pickers.

- **Command-line mode** (for automation or batch use)

- **Three conversion modes**  
  1. **In-place conversion** (_HEIC → JPEG next to originals_)

  2. **Convert → destination folder** (_HEIC → JPEG into a specified output folder (preserves structure)_)

  3. **Copy everything + convert HEIC** (_Copies all files AND converts HEIC → JPEG_)

- **Advanced options**  
  - Adjustable JPEG quality (`--quality`)  
  - Dry-run mode (`--dry-run`)  
  - Delete originals after successful conversion (`--delete-originals`)  
  - Configurable concurrency (`--workers`)

- **Cross-platform**  
  Uses `github.com/ncruces/zenity` for native dialogs.

---

## 🚀 Getting Started

### Requirements

- Go 1.20+ (for compiling)  
- Or download a prebuilt binary (in future releases)

---

## 📦 Installation

### Clone & build manually

```sh
git clone https://github.com/yourname/heic2jpeg
cd heic2jpeg
make build
```

The compiled binary will be in:

```sh
target/heic2jpeg
```

(On Windows: target/heic2jpeg.exe)

## 🎉 Guided Mode (default)

Just run:

```sh
./heic2jpeg
```

Guided mode will:

1. Ask for the conversion type

2. Offer advanced settings (quality, dry-run, delete originals)

3. Open a native source folder picker

4. Show how many files and HEIC photos were found

5. Ask for a destination folder if required

6. Display a live progress bar during conversion

Perfect for beginners or non-technical users.

## 🔧 CLI Mode

To use the command-line interface:

```sh
./heic2jpeg cli [flags]
```

### Modes (choose ONE)

| Mode                   | Flag        | Description                                                          |
| ---------------------- | ----------- | -------------------------------------------------------------------- |
| In-place               | `--inplace` | Converts HEIC → JPEG next to the originals                           |
| Convert to destination | `--convert` | Converts HEIC → JPEG into a destination folder (preserves structure) |
| Copy + Convert         | `--copy`    | Copies all files AND converts HEIC → JPEG                            |

### Required Flags

| Flag                 | Description                                           |
| -------------------- | ----------------------------------------------------- |
| `--source PATH`      | Source folder                                         |
| `--destination PATH` | Destination folder (only for `--convert` or `--copy`) |

### Optional Flags

| Flag                 | Default   | Description                                             |
| -------------------- | --------- | ------------------------------------------------------- |
| `--quality N`        | 90        | JPEG quality (1–100)                                    |
| `--dry-run`          | false     | Simulates conversion without writing files              |
| `--delete-originals` | false     | Deletes original HEIC files after successful conversion |
| `--workers N`        | CPU cores | Number of parallel workers                              |

### 📘 CLI Examples

**In-place conversion:**

```sh
./heic2jpeg cli --inplace --source "/photos"
```

**Convert to another folder (preserve structure):**

```sh
./heic2jpeg cli --convert --source "/input" --destination "/output"
```

**Copy everything + convert HEIC:**

```sh
./heic2jpeg cli --copy --source "/iphone-dump" --destination "/ready-for-desktop"
```

**Dry-run (no files written):**

```sh
./heic2jpeg cli --convert --source "/photos" --destination "/export" --dry-run
```

**Custom quality + delete originals:**

```sh
./heic2jpeg cli --inplace --source "/photos" --quality 75 --delete-originals
```

## 📂 Folder Structure

```
.
├── cmd/
│   └── heic2jpeg/      # main.go entrypoint
├── internal/
│   └── app/            # application logic
├── target/             # build output from Makefile
└── Makefile
```

## 🧪 Running Tests

Run all tests:

```sh
make test
```

Or directly:

```sh
go test ./...
```

## 🤝 Contributing

Pull requests and issues are welcome! Ideas for improvements—UI/UX, performance, or features—are appreciated.

## 📄 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

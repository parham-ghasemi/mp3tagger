# MP3Tagger

`mp3tagger` is a lightweight command-line tool written in Go that automatically scrapes album tracklists from Genius.com and updates the ID3 tags (`Title`, `Artist`, `Album`, and `Track Number`) of your local MP3 files.

It was born out of a personal frustration: downloading albums via tools like `yt-dlp` often leaves you with messy filenames cluttered with repetitive tokens (e.g., `[Official Audio]`, `- Single`, or the artist's name duplicated on every track) and completely missing metadata.

`mp3tagger` solves this by stripping out common noise tokens globally, calculating a weighted string-similarity score, and mapping files to the correct official tracklist seamlessly.

---

## Features

- **Smart Scraping:** Pulls official track numbers, titles, primary artists, and album names directly from any Genius.com album URL.
- **Global Token-Frequency Filtering:** Analyzes your entire folder batch. If a specific word (token) appears in more than a certain percentage of your files (e.g., `60%`), it’s flagged as noise and ignored during comparison.
- **Weighted Similarity Algorithm:** Lowers the matching weight of generic modifiers like `ft.`, `feat.`, `official`, or `audio` to prevent false positives.
- **Dry-Run Mode (`-n`):** Preview exactly what matching logic will execute and see the confidence scores before writing a single byte to your files.
- **Fault Resilient:** If a single MP3 file in a massive batch is locked or corrupted, the tool logs the error and gracefully moves on to the next track rather than crashing.

---

## Installation & Quick Start

No Go installation or source code compilation required! We provide pre-compiled, standalone binaries for every major operating system.

1. Head over to the [Releases](https://github.com/parham-ghasemi/mp3tagger/releases) page and download the executable file that matches your system architecture.
2. Open your terminal or command prompt and navigate to the folder where you downloaded the file.
3. Run the application directly by supplying your target Genius URL and the directory containing your MP3 files:

### Windows

```cmd
.\mp3tagger-windows.exe -u <genius-album-url> -d <path-to-mp3-folder>
```

### macOS (Apple Silicon / M-Series)

```Bash
# Give the binary permission to run, then execute
chmod +x mp3tagger-mac-arm64
./mp3tagger-mac-arm64 -u <genius-album-url> -d <path-to-mp3-folder>
```

### macOS (Intel Chips)

```Bash
chmod +x mp3tagger-mac-intel
./mp3tagger-mac-intel -u <genius-album-url> -d <path-to-mp3-folder>
```

### Linux (Universal)

```Bash
chmod +x mp3tagger-linux
./mp3tagger-linux -u <genius-album-url> -d <path-to-mp3-folder>
```

**Note for Developers:** If you prefer to run or build the project directly from source code, ensure you have Go 1.22+ installed and run `go run main.go [flags]` from the repository root.

| **Long Flag**       | **Short Flag** | **Default** | **Description**                                                                       |
| ------------------- | -------------- | ----------- | ------------------------------------------------------------------------------------- |
| `--genius-url`      | `-u`           | _Required_  | The Genius.com album URL containing the tracklist information.                        |
| `--directory`       | `-d`           | _Required_  | Path to the local directory where your target MP3 files are stored.                   |
| `--ignore-freq-pct` | `-p`           | `60`        | Threshold %; tokens appearing in more than this % of files are filtered out as noise. |
| `--min-match-score` | `-s`           | `20`        | Minimum confidence score (0-100) required to automatically accept a track match.      |
| `--dry-run`         | `-n`           | `false`     | Preview matching logic and print summary tables without modifying file tags.          |

---

## How It Works

The tool processes your audio library using a modular architecture split into distinct domain packages:

1. **`scraper`**: Leverages `goquery` to cleanly pull and parse track arrays directly from the Genius.com DOM layout.
2. **`cleaner`**: Preps your filenames for matching. It strips out trailing YouTube video IDs (like `[dQw4w9WgXcQ]`) and filters out repetitive text clutter across your files by building a global token-frequency map.
3. **`compare`**: Normalizes characters and applies an intersection-over-union string matching algorithm. Tokens are passed through an internal lookup weight matrix so background terms like `official` or `pt` don't throw off the similarity indexing.
4. **`tagger`**: Directly manipulates the underlying audio file frames using the `id3v2` binary parser, executing fast UTF-8 textual payload injection.

---

## Dependencies

This tool leans on these libraries:

- [PuerkitoBio/goquery](https://github.com/PuerkitoBio/goquery) - HTML parsing
- [bogem/id3v2](https://github.com/bogem/id3v2) - High-performance ID3v2 tagging

## Todo

- [ ] Add Album Cover
- [ ] Move Matching Logic out of main

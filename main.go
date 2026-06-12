package main

import (
	"flag"
	"fmt"
	"log"
	removecommontokens "mp3tagger/removeCommonTokens"
	"mp3tagger/scraper"
	"mp3tagger/tagger"
	tokencompare "mp3tagger/tokenCompare"
	"os"
	"path/filepath"
	"strings"
)

type Match struct {
	Track scraper.Track
	Score int32
}


func main() {
	var (
		geniusUrl     string
		directory     string
		ignoreFreqPct int
		minMatchScore int
		dryRun        bool
	)

	flag.StringVar(&geniusUrl, "genius-url", "", "The Genius.com album URL containing the correct tracklist information")
	flag.StringVar(&geniusUrl, "u", "", "The Genius.com album URL containing the correct tracklist information (shorthand)")

	flag.StringVar(&directory, "directory", "", "Path to the local directory where your target MP3 files are stored")
	flag.StringVar(&directory, "d", "", "Path to the local directory where your target MP3 files are stored (shorthand)")

	flag.IntVar(&ignoreFreqPct, "ignore-freq-pct", 60, "Threshold percentage (0-100); tokens appearing in more than this % of files are filtered out as too common")
	flag.IntVar(&ignoreFreqPct, "p", 60, "Threshold percentage (0-100); tokens appearing in more than this % of files are filtered out as too common (shorthand)")

	flag.IntVar(&minMatchScore, "min-match-score", 20, "Minimum confidence score (0-100) required to automatically accept a track match")
	flag.IntVar(&minMatchScore, "s", 20, "Minimum confidence score (0-100) required to automatically accept a track match (shorthand)")

	flag.BoolVar(&dryRun, "dry-run", false, "Preview matching logic and print what changes would be made without modifying file tags")
	flag.BoolVar(&dryRun, "n", false, "Preview matching logic and print what changes would be made without modifying file tags (shorthand)")

	flag.Parse()

	if geniusUrl == "" || directory == "" {
		flag.Usage()
		os.Exit(1)
	}

	fmt.Println("Starting app")

	scrapeRes, err := scraper.ScrapePage(geniusUrl)
	if err != nil {
		log.Fatalf("Failed to scrape page: %v", err)
	}
	fmt.Printf("Scrape result: %d tracks found\n", len(scrapeRes))


	entries, err := os.ReadDir(directory)
	if err != nil {
		log.Fatalf("Failed to read directory: %v", err)
	}

	fmt.Println("Searching for MP3 files in top-level folder")
	files := make([]string, 0, len(entries))
	for _, entry := range entries {
		if strings.ToLower((filepath.Ext((entry.Name())))) == ".mp3" {
			files = append(files, entry.Name())
		}
	}

	fmt.Printf("Scanned %d directory entries; found %d MP3 files to process.\n", len(entries), len(files))

	cleanFiles := removecommontokens.RemoveCommonTokens(files, ignoreFreqPct)

	for index, file := range cleanFiles {
		bestScore := 0
		var bestTrack scraper.Track

		for _, track := range scrapeRes {
			score := tokencompare.TokenCompare(file, track.Title)

			if score > bestScore {
				bestScore = score
				bestTrack = track
			}
		}

		if bestScore < minMatchScore {
			fmt.Printf("\nSKIP: %s — low confidence (score= %d, min= %d). Best match: %s\n", file, bestScore, minMatchScore, bestTrack.Title)
			continue
		}

		fmt.Printf("\nMATCH: %s -> %s (score= %d)\n", file, bestTrack.Title, bestScore)
		fileFullPath := filepath.Join(directory, files[index])
		fmt.Printf("Editing tags on %s\n", fileFullPath)
		if dryRun {
			log.Printf("Dry run: skipping tag edit for %s", fileFullPath)
		} else {
			if err := tagger.Tagger(fileFullPath, bestTrack.TrackNumber, bestTrack.Title, bestTrack.AlbumArtist, bestTrack.AlbumName); err != nil {
				log.Printf("Failed to edit tags on %s: %v\n", fileFullPath, err)
			}
			log.Printf("Tags edit success: %s", fileFullPath)
		}
	}
}


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
	var geniusUrl = flag.String("genius-url", "", "The genius url, the tags will be fetched from.")
	var directory = flag.String("directory", "", "The directory in which your MP3 files are located.")
	var ignoreFreqPct = flag.Int("ignore-freq-pct", 60, "The maximum percentage of files a token can appear in; tokens exceeding this are ignored as too common")
	var minMatchScore = flag.Int("min-match-score", 20, "The minimum confidence score required to accept a track match")
	var dryRun = flag.Bool("dry-run", false, "If set, run without editing any tags")

	flag.Parse()

	if *geniusUrl == "" || *directory == "" {
		flag.Usage()
		os.Exit(1)
	}

	fmt.Println("Starting app")

	scrapeRes, err := scraper.ScrapePage(*geniusUrl)
	if err != nil {
		log.Fatalf("Failed to scrape page: %v", err)
	}
	fmt.Printf("Scrape result: %d tracks found\n", len(scrapeRes))


	entries, err := os.ReadDir(*directory)
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

	cleanFiles := removecommontokens.RemoveCommonTokens(files, *ignoreFreqPct)

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

		if bestScore < *minMatchScore {
			fmt.Printf("\nSKIP: %s — low confidence (score= %d, min= %d). Best match: %s\n", file, bestScore, *minMatchScore, bestTrack.Title)
			continue
		}

		fmt.Printf("\nMATCH: %s -> %s (score= %d)\n", file, bestTrack.Title, bestScore)
		fileFullPath := filepath.Join(*directory, files[index])
		fmt.Printf("Editing tags on %s\n", fileFullPath)
		if *dryRun {
			log.Printf("Dry run: skipping tag edit for %s", fileFullPath)
		} else {
			if err := tagger.Tagger(fileFullPath, bestTrack.TrackNumber, bestTrack.Title, bestTrack.AlbumArtist, bestTrack.AlbumName); err != nil {
				log.Printf("Failed to edit tags on %s: %v\n", fileFullPath, err)
			}
			log.Printf("Tags edit success: %s", fileFullPath)
		}
	}
}


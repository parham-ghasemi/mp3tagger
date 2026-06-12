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

	flag.Parse()

	if *geniusUrl == "" || *directory == "" {
		flag.Usage()
		os.Exit(1)
	}

	fmt.Println("Starting app...")

	scrapeRes := scraper.ScrapePage(*geniusUrl)
	fmt.Println("Scrape Result: ")
	fmt.Println(scrapeRes)


	entries, err := os.ReadDir(*directory)
	if err != nil {
		log.Fatalf("Failed to read directory: %v", err)
	}

	fmt.Println("Searching for MP3 files in top-level folder...")
	files := []string {}
	for _, entry := range entries {
		if strings.ToLower((filepath.Ext((entry.Name())))) == ".mp3" {
			files = append(files, entry.Name())
		}
	}

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
			fmt.Println("\nSKIP (low confidence)")
			fmt.Println("------")
			fmt.Printf("Clean MP3 Name: %v \n", file)
			fmt.Printf("MATCH NAME: %v \n", bestTrack.Title)
			fmt.Printf("Score: %v\n", bestScore) 

			continue
		}

		fmt.Println("\nMATCH")
		fmt.Println("------")
		fmt.Printf("Clean MP3 Name: %v \n", file)
		fmt.Printf("MATCH NAME: %v \n", bestTrack.Title)
		fmt.Printf("Score: %v \n", bestScore) 
		fileFullPath := filepath.Join(*directory, files[index]) 
		fmt.Println("Editing tag on file ", fileFullPath) 
		tagger.Tagger(fileFullPath, bestTrack.TrackNumber, bestTrack.Title, bestTrack.AlbumArtist, bestTrack.AlbumName) 
		fmt.Println("Tags edit success.")
	}
}


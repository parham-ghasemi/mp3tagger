package main

import (
	"bufio"
	"fmt"
	"log"
	removecommontokens "mp3tagger/removeCommonTokens"
	"mp3tagger/scraper"
	"mp3tagger/tagger"
	"os"
	"path/filepath"
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

type Match struct {
	Track scraper.Track
	Score int32
}

func main() {
	fmt.Println("Starting app...")

	// Scrape the genius url
	var genius_url string
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter Genius URL:")
	scanner.Scan()
	genius_url = scanner.Text()
	if err := scanner.Err(); err != nil {
		fmt.Fprint(os.Stderr, "reading standard input: ", err)
	}

	scrapeRes := scraper.ScrapePage(genius_url)
	fmt.Println("Scrape Result: ")
	fmt.Println(scrapeRes)


	// Loop through folder
	var folder_path string
	fmt.Println("Enter Folder URL:")
	scanner.Scan()
	folder_path = scanner.Text()
	if err := scanner.Err(); err != nil {
		fmt.Fprint(os.Stderr, "reading standard input: ", err)
	}
	entries, err := os.ReadDir(folder_path)
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

	cleanFiles := removecommontokens.RemoveCommonTokens(files)

	for index, file := range cleanFiles {
		bestScore := 999999
		var bestTrack scraper.Track

		for _, track := range scrapeRes {
			score := fuzzy.LevenshteinDistance(file, strings.ToLower(track.Title))

			if score < bestScore {
				bestScore = score
				bestTrack = track
			}
		}

		maxScore := 40

		if bestScore > maxScore {
			fmt.Println("\nSKIP (low confidence)")
			fmt.Println("------")
			fmt.Printf("Clean MP3 Name: %v \n", file)
			fmt.Printf("MATCH NAME: %v \n", bestTrack.Title)
			fmt.Printf("Score: %v \n", bestScore) 

			continue
		}

		fmt.Println("\nMATCH")
		fmt.Println("------")
		fmt.Printf("Clean MP3 Name: %v \n", file)
		fmt.Printf("MATCH NAME: %v \n", bestTrack.Title)
		fmt.Printf("Score: %v \n", bestScore) 
		fileFullPath := filepath.Join(folder_path, files[index]) 
		fmt.Println("Editing tag on file ", fileFullPath) 
		tagger.Tagger(fileFullPath, bestTrack.TrackNumber, bestTrack.Title, bestTrack.AlbumArtist, bestTrack.AlbumName) 
		fmt.Println("Tags edit success.")
	}
}


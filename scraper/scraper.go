package scraper

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Track struct {
	TrackNumber string
	Title string
	AlbumName string
	AlbumArtist string
}

// ScrapePage starts with a Capital letter so main.go can use it
func ScrapePage(url string) []Track {
	var output []Track

	fmt.Println("Scraping: ", url, "...")

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	AlbumTitle := strings.TrimSpace(doc.Find(".header_with_cover_art-primary_info-title").First().Text())
	AlbumPrimaryArtists := strings.TrimSpace(doc.Find(".header_with_cover_art-primary_info-primary_artist").First().Text())

	doc.Find(".chart_row-content-title").Each(func(i int, s *goquery.Selection) {
		clonedSelection := s.Clone()
		clonedSelection.Find("span").Remove()
		text := clonedSelection.Text()
		cleanText := strings.TrimSpace(text)
		result := Track{TrackNumber: strconv.Itoa(i + 1), Title: cleanText, AlbumName: AlbumTitle, AlbumArtist: AlbumPrimaryArtists}
		

		output = append(output, result)
	})

	return output
}

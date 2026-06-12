package tokencompare

/*
	[ ] Weighted tokens
	[X] Normalize strings before matching
*/

import (
	"slices"
	"strings"
)

func Normalize(s string) string{
	s = strings.ToLower((s))

	replacer := strings.NewReplacer( ".", "", ",", "", "(", "", ")", "", "-", " ",)

	return replacer.Replace(s)
}

func TokenCompare(s1, s2 string) int {
	s1t := strings.Fields(Normalize(s1))
	s2t := strings.Fields(Normalize(s2))

	common := map[string]bool{}

	for _, token := range s1t {
		if slices.Contains(s2t, token) {
			common[token] = true
		}
	}

	similarity := int(float64(len(common)) / (float64(len(s1t)) + float64(len(s2t)) - float64(len(common))) * 100)
	return similarity
}